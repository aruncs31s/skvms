package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aruncs31s/skvms/internal/logger"
	"go.uber.org/zap"
)

const (
	// RepoURL is the GitHub repository containing the ESP32 firmware source.
	RepoURL = "https://github.com/aruncs31s/Kannur-Solar-Battery-Monitoring-System-Microcontroller-Codes.git"

	// RepoName is the directory name for the cloned repo.
	RepoName = "esp32-firmware-source"
)

// CloneOrPullRepo ensures the firmware source repo exists at baseDir/RepoName.
// If it doesn't exist, it clones it. If it does, it pulls the latest changes.
func CloneOrPullRepo(baseDir string) (string, error) {
	repoDir := filepath.Join(baseDir, RepoName)

	if _, err := os.Stat(filepath.Join(repoDir, ".git")); os.IsNotExist(err) {
		// Clone the repository
		logger.GetLogger().Info("Cloning firmware source repository",
			zap.String("url", RepoURL),
			zap.String("target", repoDir),
		)

		if err := os.MkdirAll(baseDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create base directory: %w", err)
		}

		cmd := exec.Command("git", "clone", "--depth", "1", RepoURL, repoDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("git clone failed: %w\nOutput: %s", err, string(output))
		}

		logger.GetLogger().Info("Repository cloned successfully",
			zap.String("path", repoDir),
		)
	} else {
		// Pull latest changes
		logger.GetLogger().Info("Pulling latest firmware source",
			zap.String("path", repoDir),
		)

		cmd := exec.Command("git", "-C", repoDir, "pull", "--ff-only")
		output, err := cmd.CombinedOutput()
		if err != nil {
			logger.GetLogger().Warn("Git pull failed, using existing source",
				zap.String("output", string(output)),
				zap.Error(err),
			)
		}
	}

	return repoDir, nil
}

// CopyRepoForBuild creates an isolated copy of the source repo for a specific build.
// This prevents concurrent builds from interfering with each other.
func CopyRepoForBuild(sourceDir, buildID string) (string, error) {
	buildDir := filepath.Join(filepath.Dir(sourceDir), "builds", buildID)

	if err := os.MkdirAll(filepath.Dir(buildDir), 0755); err != nil {
		return "", fmt.Errorf("failed to create builds directory: %w", err)
	}

	cmd := exec.Command("cp", "-a", sourceDir, buildDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to copy repo for build: %w\nOutput: %s", err, string(output))
	}

	logger.GetLogger().Info("Created build copy",
		zap.String("build_id", buildID),
		zap.String("build_dir", buildDir),
	)

	return buildDir, nil
}

// CleanupBuild removes the build directory after use.
func CleanupBuild(buildDir string) {
	if err := os.RemoveAll(buildDir); err != nil {
		logger.GetLogger().Warn("Failed to cleanup build directory",
			zap.String("dir", buildDir),
			zap.Error(err),
		)
	}
}
