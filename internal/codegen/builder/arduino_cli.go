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

const (
	// Default FQBN for ESP32 DevKit V1
	defaultESP32FQBN = "esp32:esp32:esp32"

	// ESP8266 NodeMCU V2
	defaultESP8266FQBN = "esp8266:esp8266:nodemcuv2"
)

// ArduinoCLIBuilder implements BuildStrategy using Arduino CLI.
type ArduinoCLIBuilder struct {
	FQBN string // Fully Qualified Board Name (e.g., "esp32:esp32:esp32")
}

func NewArduinoCLIBuilder(fqbn string) *ArduinoCLIBuilder {
	if fqbn == "" {
		fqbn = defaultESP32FQBN
	}
	return &ArduinoCLIBuilder{FQBN: fqbn}
}

func (a *ArduinoCLIBuilder) Name() string {
	return "Arduino CLI"
}

func (a *ArduinoCLIBuilder) IsAvailable() bool {
	_, err := exec.LookPath("arduino-cli")
	return err == nil
}

func (a *ArduinoCLIBuilder) Build(ctx context.Context, projectDir string) (*BuildResult, error) {
	logger.GetLogger().Info("Building firmware with Arduino CLI",
		zap.String("project_dir", projectDir),
		zap.String("fqbn", a.FQBN),
	)

	// Ensure the board core is installed
	if err := a.ensureCoreInstalled(ctx); err != nil {
		logger.GetLogger().Warn("Failed to install board core (may already be installed)",
			zap.Error(err),
		)
	}

	// Install library dependencies
	if err := a.installLibraries(ctx); err != nil {
		logger.GetLogger().Warn("Failed to install libraries (may already be installed)",
			zap.Error(err),
		)
	}

	// For Arduino CLI, the sketch must be in a folder with the same name as the .ino file.
	// PlatformIO projects use src/main.cpp. We need to create a wrapper .ino file.
	sketchDir, err := a.prepareSketchDir(projectDir)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare sketch directory: %w", err)
	}

	outputDir := filepath.Join(projectDir, "build", "arduino-cli-output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	cmd := exec.CommandContext(ctx,
		"arduino-cli", "compile",
		"--fqbn", a.FQBN,
		"--output-dir", outputDir,
		"--libraries", filepath.Join(projectDir, "lib"),
		sketchDir,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.GetLogger().Error("Arduino CLI build failed",
			zap.String("output", string(output)),
			zap.Error(err),
		)
		return nil, fmt.Errorf("arduino-cli build failed: %w\nOutput: %s", err, string(output))
	}

	logger.GetLogger().Info("Arduino CLI build succeeded",
		zap.String("output_tail", tailString(string(output), 500)),
	)

	// Find the compiled binary
	binaryPath, err := a.findBinary(outputDir)
	if err != nil {
		return nil, fmt.Errorf("build succeeded but binary not found: %w", err)
	}

	info, err := os.Stat(binaryPath)
	if err != nil {
		return nil, fmt.Errorf("cannot stat binary: %w", err)
	}

	return &BuildResult{
		BinaryPath: binaryPath,
		BoardFQBN:  a.FQBN,
		Size:       info.Size(),
	}, nil
}

func (a *ArduinoCLIBuilder) ensureCoreInstalled(ctx context.Context) error {
	parts := strings.SplitN(a.FQBN, ":", 3)
	if len(parts) < 2 {
		return fmt.Errorf("invalid FQBN: %s", a.FQBN)
	}
	core := parts[0] + ":" + parts[1]

	// Add board manager URL for ESP boards
	var boardURL string
	switch parts[0] {
	case "esp32":
		boardURL = "https://raw.githubusercontent.com/espressif/arduino-esp32/gh-pages/package_esp32_index.json"
	case "esp8266":
		boardURL = "https://arduino.esp8266.com/stable/package_esp8266com_index.json"
	}

	if boardURL != "" {
		cmd := exec.CommandContext(ctx, "arduino-cli", "config", "add", "board_manager.additional_urls", boardURL)
		_ = cmd.Run() // Ignore errors if already added

		cmd = exec.CommandContext(ctx, "arduino-cli", "core", "update-index")
		_ = cmd.Run()
	}

	cmd := exec.CommandContext(ctx, "arduino-cli", "core", "install", core)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("core install failed: %w\nOutput: %s", err, string(output))
	}
	return nil
}

func (a *ArduinoCLIBuilder) installLibraries(ctx context.Context) error {
	libs := []string{"ArduinoJson"}
	for _, lib := range libs {
		cmd := exec.CommandContext(ctx, "arduino-cli", "lib", "install", lib)
		_ = cmd.Run() // Best effort
	}
	return nil
}

func (a *ArduinoCLIBuilder) prepareSketchDir(projectDir string) (string, error) {
	// If there's already a .ino file, use the project dir directly
	entries, err := os.ReadDir(projectDir)
	if err != nil {
		return "", err
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".ino") {
			return projectDir, nil
		}
	}

	// Create a sketch wrapper that includes main.cpp
	sketchDir := filepath.Join(projectDir, "arduino_sketch")
	if err := os.MkdirAll(sketchDir, 0755); err != nil {
		return "", err
	}

	// Create .ino file that includes the main.cpp
	inoContent := fmt.Sprintf(`// Auto-generated Arduino sketch wrapper
// This file includes the PlatformIO main.cpp for Arduino CLI compatibility
#include "%s"
`, filepath.Join(projectDir, "src", "main.cpp"))

	inoPath := filepath.Join(sketchDir, "arduino_sketch.ino")
	if err := os.WriteFile(inoPath, []byte(inoContent), 0644); err != nil {
		return "", err
	}

	return sketchDir, nil
}

func (a *ArduinoCLIBuilder) findBinary(outputDir string) (string, error) {
	// Arduino CLI outputs files like sketch.ino.bin or sketch.ino.esp32.bin
	var binPath string

	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".bin") {
			binPath = path
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error searching for binary: %w", err)
	}
	if binPath == "" {
		return "", fmt.Errorf("no .bin file found in %s", outputDir)
	}
	return binPath, nil
}

func (a *ArduinoCLIBuilder) Upload(ctx context.Context, projectDir string, deviceIP string) error {
	logger.GetLogger().Info("Uploading firmware via Arduino CLI OTA",
		zap.String("project_dir", projectDir),
		zap.String("device_ip", deviceIP),
		zap.String("fqbn", a.FQBN),
	)

	// Arduino CLI doesn't have native OTA support like PlatformIO.
	// We use espota.py which comes with the ESP32 Arduino core.
	// Alternatively, build first then use espota.

	// First, build to get the binary
	result, err := a.Build(ctx, projectDir)
	if err != nil {
		return fmt.Errorf("build failed before upload: %w", err)
	}

	// Try espota.py for OTA upload
	espotaPath, err := findEspota()
	if err != nil {
		return fmt.Errorf("espota.py not found, cannot perform OTA upload: %w", err)
	}

	cmd := exec.CommandContext(ctx,
		"python3", espotaPath,
		"-i", deviceIP,
		"-f", result.BinaryPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.GetLogger().Error("Arduino CLI OTA upload failed",
			zap.String("output", string(output)),
			zap.Error(err),
		)
		return fmt.Errorf("OTA upload failed: %w\nOutput: %s", err, string(output))
	}

	logger.GetLogger().Info("Arduino CLI OTA upload succeeded",
		zap.String("device_ip", deviceIP),
	)
	return nil
}

// findEspota looks for the espota.py script in common locations.
func findEspota() (string, error) {
	commonPaths := []string{
		filepath.Join(os.Getenv("HOME"), ".arduino15", "packages", "esp32", "hardware", "esp32"),
		filepath.Join(os.Getenv("HOME"), ".platformio", "packages", "framework-arduinoespressif32", "tools"),
		"/usr/share/arduino/hardware/espressif/esp32/tools",
	}

	for _, basePath := range commonPaths {
		var found string
		_ = filepath.Walk(basePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.Name() == "espota.py" {
				found = path
				return filepath.SkipAll
			}
			return nil
		})
		if found != "" {
			return found, nil
		}
	}

	// Try in PATH
	path, err := exec.LookPath("espota.py")
	if err == nil {
		return path, nil
	}

	return "", fmt.Errorf("espota.py not found in common locations")
}
