package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aruncs31s/skvms/internal/logger"
	"go.uber.org/zap"
)

// PlatformIOBuilder implements BuildStrategy using PlatformIO CLI (pio).
type PlatformIOBuilder struct{}

func NewPlatformIOBuilder() *PlatformIOBuilder {
	return &PlatformIOBuilder{}
}

func (p *PlatformIOBuilder) Name() string {
	return "PlatformIO"
}

func (p *PlatformIOBuilder) IsAvailable() bool {
	_, err := exec.LookPath("pio")
	if err != nil {
		// Also try "platformio"
		_, err = exec.LookPath("platformio")
	}
	return err == nil
}

func (p *PlatformIOBuilder) pioBin() string {
	path, err := exec.LookPath("pio")
	if err == nil {
		return path
	}
	path, err = exec.LookPath("platformio")
	if err == nil {
		return path
	}
	return "pio"
}

func (p *PlatformIOBuilder) Build(ctx context.Context, projectDir string) (*BuildResult, error) {
	logger.GetLogger().Info("Building firmware with PlatformIO",
		zap.String("project_dir", projectDir),
	)

	cmd := exec.CommandContext(ctx, p.pioBin(), "run", "-d", projectDir)
	cmd.Dir = projectDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.GetLogger().Error("PlatformIO build failed",
			zap.String("output", string(output)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("platformio build failed: %w\nOutput: %s", err, string(output))
	}

	logger.GetLogger().Info("PlatformIO build succeeded",
		zap.String("output_tail", tailString(string(output), 500)),
	)

	// Find the compiled binary. PlatformIO puts it in .pio/build/<env>/firmware.bin
	binaryPath, err := p.findBinary(projectDir)
	if err != nil {
		return nil, fmt.Errorf("build succeeded but binary not found: %w", err)
	}

	info, err := os.Stat(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("cannot stat binary: %w", err)
	}

	return &BuildResult{
		BinaryPath: binaryPath,
		Size:       info.Size(),
	}, nil
}

func (p *PlatformIOBuilder) findBinary(projectDir string) (string, error) {
	pioBuildDir := filepath.Join(projectDir, ".pio", "build")

	entries, err := os.ReadDir(pioBuildDir)
	if err != nil {
		return "", fmt.Errorf("cannot read .pio/build directory: %w", err)
	}

	// Look through each env directory for firmware.bin
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		firmwarePath := filepath.Join(pioBuildDir, entry.Name(), "firmware.bin")
		if _, err := os.Stat(firmwarePath); err == nil {
			return firmwarePath, nil
		}
	}

	return "", fmt.Errorf("no firmware.bin found in %s", pioBuildDir)
}

func (p *PlatformIOBuilder) Upload(ctx context.Context, projectDir string, deviceIP string) error {
	logger.GetLogger().Info("Uploading firmware via PlatformIO OTA",
		zap.String("project_dir", projectDir),
		zap.String("device_ip", deviceIP),
	)

	cmd := exec.CommandContext(ctx,
		p.pioBin(), "run", "-d", projectDir,
		"--target", "upload",
		"--upload-port", deviceIP,
	)
	cmd.Dir = projectDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.GetLogger().Error("PlatformIO upload failed",
			zap.String("output", string(output)),
			zap.Error(err),
		)
		return fmt.Errorf("platformio OTA upload failed: %w\nOutput: %s", err, string(output))
	}

	logger.GetLogger().Info("PlatformIO OTA upload succeeded",
		zap.String("device_ip", deviceIP),
	)
	return nil
}

func tailString(s string, maxLen int) string {
	lines := strings.Split(s, "\n")
	result := strings.Join(lines, "\n")
	if len(result) > maxLen {
		return result[len(result)-maxLen:]
	}
	return result
}
