package builder

import "context"

// BuildResult holds the result of a firmware build.
type BuildResult struct {
	BinaryPath string // Absolute path to the compiled binary (.bin)
	BoardFQBN  string // Board FQBN used for the build
	Size       int64  // Binary file size in bytes
}

// BuildStrategy defines the interface for building ESP32 firmware.
// Implementations can use PlatformIO or Arduino CLI.
type BuildStrategy interface {
	// Name returns the display name of the build tool (e.g., "PlatformIO", "Arduino CLI").
	Name() string

	// IsAvailable checks whether the build tool is installed and accessible.
	IsAvailable() bool

	// Build compiles the project at projectDir and returns the path to the binary.
	Build(ctx context.Context, projectDir string) (*BuildResult, error)

	// Upload flashes the firmware binary to the ESP32 at the given IP via OTA.
	Upload(ctx context.Context, projectDir string, deviceIP string) error
}
