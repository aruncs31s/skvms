package builder

import (
	"fmt"

	"github.com/aruncs31s/skvms/internal/logger"
	"go.uber.org/zap"
)

// Resolve returns the best available BuildStrategy.
// It prefers PlatformIO over Arduino CLI if both are available.
// If a preferred tool is specified, it tries that first.
func Resolve(preferred string) (BuildStrategy, error) {
	strategies := map[string]BuildStrategy{
		"platformio":  NewPlatformIOBuilder(),
		"arduino-cli": NewArduinoCLIBuilder(""),
	}

	// If user specified a preference, try that first
	if preferred != "" {
		if s, ok := strategies[preferred]; ok {
			if s.IsAvailable() {
				logger.GetLogger().Info("Using preferred build strategy",
					zap.String("strategy", s.Name()),
				)
				return s, nil
			}
			logger.GetLogger().Warn("Preferred build tool not available, trying alternatives",
				zap.String("preferred", preferred),
			)
		}
	}

	// Auto-detect: prefer PlatformIO, fall back to Arduino CLI
	order := []string{"platformio", "arduino-cli"}
	for _, key := range order {
		s := strategies[key]
		if s.IsAvailable() {
			logger.GetLogger().Info("Auto-detected build strategy",
				zap.String("strategy", s.Name()),
			)
			return s, nil
		}
	}

	return nil, fmt.Errorf("no build tool available: install PlatformIO CLI or Arduino CLI")
}

// ListAvailable returns which build strategies are currently installed.
func ListAvailable() []string {
	available := []string{}
	strategies := []BuildStrategy{
		NewPlatformIOBuilder(),
		NewArduinoCLIBuilder(""),
	}
	for _, s := range strategies {
		if s.IsAvailable() {
			available = append(available, s.Name())
		}
	}
	return available
}
