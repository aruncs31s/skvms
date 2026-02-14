package codegen

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aruncs31s/skvms/internal/codegen/builder"
	"github.com/aruncs31s/skvms/internal/codegen/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"go.uber.org/zap"
)

// Service orchestrates the ESP32 firmware code generation pipeline:
// clone repo → replace config → build → return binary.
type Service struct {
	// workDir is the base directory for repo clones and builds.
	workDir string
}

// NewService creates a new codegen Service.
// workDir is the base directory where the firmware source will be stored.
func NewService(workDir string) *Service {
	if workDir == "" {
		workDir = filepath.Join(os.TempDir(), "skvms-codegen")
	}
	return &Service{workDir: workDir}
}

// GenerateResult holds the output of a successful firmware generation.
type GenerateResult struct {
	BuildID    string
	BinaryPath string
	BinarySize int64
	BuildTool  string
}

// Generate executes the full codegen pipeline:
//  1. Ensure the firmware source repo is cloned/up-to-date
//  2. Create an isolated copy for this build
//  3. Replace config values with the request parameters
//  4. Build the firmware using the selected strategy
//  5. Return the path to the compiled binary
func (s *Service) Generate(ctx context.Context, req dto.CodeGenRequest) (*GenerateResult, error) {
	logger.GetLogger().Info("Starting codegen pipeline",
		zap.String("device_ip", req.IP),
		zap.String("host_ip", req.HostIP),
		zap.String("wifi_ssid", req.HOSTSSID),
		zap.String("build_tool", req.BuildTool),
	)

	// Step 1: Clone or pull the source repo
	sourceDir, err := CloneOrPullRepo(s.workDir)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare source repository: %w", err)
	}

	// Step 2: Create an isolated build copy
	buildID := generateBuildID()
	buildDir, err := CopyRepoForBuild(sourceDir, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to create build copy: %w", err)
	}

	// Step 3: Replace config values
	if err := ReplaceConfig(buildDir, req); err != nil {
		CleanupBuild(buildDir)
		return nil, fmt.Errorf("failed to replace config: %w", err)
	}

	// Step 4: Resolve the build strategy
	strategy, err := builder.Resolve(req.BuildTool)
	if err != nil {
		CleanupBuild(buildDir)
		return nil, fmt.Errorf("no build tool available: %w", err)
	}

	// Step 5: Build the firmware
	result, err := strategy.Build(ctx, buildDir)
	if err != nil {
		CleanupBuild(buildDir)
		return nil, fmt.Errorf("firmware build failed: %w", err)
	}

	logger.GetLogger().Info("Codegen pipeline completed successfully",
		zap.String("build_id", buildID),
		zap.String("binary_path", result.BinaryPath),
		zap.Int64("binary_size", result.Size),
		zap.String("build_tool", strategy.Name()),
	)

	return &GenerateResult{
		BuildID:    buildID,
		BinaryPath: result.BinaryPath,
		BinarySize: result.Size,
		BuildTool:  strategy.Name(),
	}, nil
}

// Upload compiles (if needed) and flashes firmware to the ESP32 via OTA.
func (s *Service) Upload(ctx context.Context, req dto.CodeGenRequest, deviceIP string) error {
	// First generate the firmware
	result, err := s.Generate(ctx, req)
	if err != nil {
		return fmt.Errorf("firmware generation failed: %w", err)
	}
	defer s.CleanupBuild(result.BuildID)

	// Resolve build strategy
	strategy, err := builder.Resolve(req.BuildTool)
	if err != nil {
		return err
	}

	// Get the build directory from the binary path
	buildDir := filepath.Dir(filepath.Dir(result.BinaryPath))

	// Upload via OTA
	return strategy.Upload(ctx, buildDir, deviceIP)
}

// GetBinaryPath returns the path to a previously built binary.
func (s *Service) GetBinaryPath(buildID string) (string, error) {
	buildDir := filepath.Join(s.workDir, "builds", buildID)
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		return "", fmt.Errorf("build %s not found", buildID)
	}

	// Search for .bin file in the build directory
	var binPath string
	err := filepath.Walk(buildDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".bin" {
			binPath = path
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil || binPath == "" {
		return "", fmt.Errorf("binary not found for build %s", buildID)
	}

	return binPath, nil
}

// CleanupBuild removes the build directory for the given build ID.
func (s *Service) CleanupBuild(buildID string) {
	buildDir := filepath.Join(s.workDir, "builds", buildID)
	CleanupBuild(buildDir)
}

// ListAvailableTools returns the names of available build tools.
func (s *Service) ListAvailableTools() []string {
	return builder.ListAvailable()
}

// generateBuildID creates a unique build identifier using timestamp and random suffix.
func generateBuildID() string {
	return fmt.Sprintf("%d", os.Getpid()) + "-" + randomHex(8)
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = fmt.Fprintf(os.Stderr, "") // side-effect free
	f, err := os.Open("/dev/urandom")
	if err != nil {
		// Fallback to time-based
		return fmt.Sprintf("%x", os.Getpid())
	}
	defer f.Close()
	_, _ = f.Read(b)
	return fmt.Sprintf("%x", b)
}
