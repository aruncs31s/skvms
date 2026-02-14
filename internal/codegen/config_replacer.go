package codegen

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aruncs31s/skvms/internal/codegen/dto"
	"github.com/aruncs31s/skvms/internal/logger"
	"go.uber.org/zap"
)

// ReplaceConfig modifies the config.h file in the project directory with the
// values from the CodeGenRequest. This replaces WiFi credentials, backend
// server details, and the device token.
func ReplaceConfig(projectDir string, req dto.CodeGenRequest) error {
	configPath := filepath.Join(projectDir, "include", "config.h")

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config.h: %w", err)
	}

	original := string(content)
	modified := original

	// 1. Replace BACKEND_HOST
	if req.HostIP != "" {
		modified = replaceDefine(modified, "BACKEND_HOST", fmt.Sprintf(`"%s"`, req.HostIP))
	}

	// 2. Replace BACKEND_PORT
	if req.Port > 0 {
		modified = replaceDefine(modified, "BACKEND_PORT", fmt.Sprintf("%d", req.Port))
	}

	// 3. Replace TOKEN
	if req.Token != "" {
		modified = replaceDefine(modified, "TOKEN", fmt.Sprintf(`"%s"`, req.Token))
	}

	// 4. Replace WiFi SSID — replace all WIFI_SSID definitions
	if req.HOSTSSID != "" {
		modified = replaceAllDefine(modified, "WIFI_SSID", fmt.Sprintf(`"%s"`, req.HOSTSSID))
	}

	// 5. Replace WiFi PASSWORD — replace all WIFI_PASSWORD definitions
	if req.HOSTPASS != "" {
		modified = replaceAllDefine(modified, "WIFI_PASSWORD", fmt.Sprintf(`"%s"`, req.HOSTPASS))
	}

	// 6. Replace STATIC_IP_ADDRESS if device IP is provided
	if req.IP != "" {
		ipOctets := ipToOctets(req.IP)
		if ipOctets != "" {
			modified = replaceAllDefine(modified, "STATIC_IP_ADDRESS", ipOctets)
			modified = uncommentDefine(modified, "STATIC_IP")
		}
	}

	// 7. Replace DEVICE_NAME if provided
	if req.DeviceName != "" {
		modified = replaceDefine(modified, "DEVICE_NAME", fmt.Sprintf(`"%s"`, req.DeviceName))
	}

	// 8. Ensure USE_GO_BACKEND is defined
	modified = uncommentDefine(modified, "USE_GO_BACKEND")

	if modified == original {
		logger.GetLogger().Warn("No changes were made to config.h")
	}

	if err := os.WriteFile(configPath, []byte(modified), 0644); err != nil {
		return fmt.Errorf("failed to write config.h: %w", err)
	}

	logger.GetLogger().Info("Config replaced successfully",
		zap.String("config_path", configPath),
		zap.String("backend_host", req.HostIP),
		zap.Int("backend_port", req.Port),
		zap.String("wifi_ssid", req.HOSTSSID),
		zap.String("device_ip", req.IP),
	)

	return nil
}

// replaceDefine replaces the first occurrence of #define KEY <old_value> with #define KEY <new_value>.
func replaceDefine(content, key, newValue string) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?m)^(\s*#define\s+%s\s+)(.+)$`, regexp.QuoteMeta(key)))
	if pattern.MatchString(content) {
		return pattern.ReplaceAllString(content, "${1}"+newValue)
	}
	return content
}

// replaceAllDefine replaces ALL occurrences of #define KEY <old_value> with #define KEY <new_value>.
func replaceAllDefine(content, key, newValue string) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?m)^(\s*#define\s+%s\s+)(.+)$`, regexp.QuoteMeta(key)))
	return pattern.ReplaceAllString(content, "${1}"+newValue)
}

// uncommentDefine ensures a #define is active (uncomments it if needed).
func uncommentDefine(content, key string) string {
	pattern := regexp.MustCompile(fmt.Sprintf(`(?m)^(\s*)//\s*(#define\s+%s.*)$`, regexp.QuoteMeta(key)))
	return pattern.ReplaceAllString(content, "${1}${2}")
}

// ipToOctets converts "192.168.1.50" to "192, 168, 1, 50" for C macro format.
func ipToOctets(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return ""
	}
	return strings.Join(parts, ", ")
}
