package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var requiredExecutables = []string{"yt-dlp", "ffmpeg"}

type Config struct {
	RequiredExecutables map[string]string
}

func (c *Config) MustExecute() (*Config, error) {
	var missing []string

	if c.RequiredExecutables == nil {
		c.RequiredExecutables = make(map[string]string)
	}

	for _, executable := range requiredExecutables {
		if path, err := resolveToolPath(executable); err != nil {
			missing = append(missing, executable)
		} else {
			c.RequiredExecutables[executable] = path
		}
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("startup dependency check failed: missing %v", missing)
	}

	return c, nil
}

func resolveToolPath(toolName string) (string, error) {
	if path, err := exec.LookPath(toolName); err == nil {
		return path, nil
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	var patterns []string

	switch toolName {
	case "yt-dlp":
		patterns = []string{
			filepath.Join(localAppData, "Microsoft", "WinGet", "Links", "yt-dlp.exe"),
			filepath.Join(localAppData, "Microsoft", "WinGet", "Packages", "yt-dlp.yt-dlp_*", "yt-dlp.exe"),
		}
	case "ffmpeg":
		patterns = []string{
			filepath.Join(localAppData, "Microsoft", "WinGet", "Links", "ffmpeg.exe"),
			filepath.Join(localAppData, "Microsoft", "WinGet", "Packages", "Gyan.FFmpeg_*", "*", "bin", "ffmpeg.exe"),
			filepath.Join(localAppData, "Microsoft", "WinGet", "Packages", "yt-dlp.FFmpeg_*", "*", "bin", "ffmpeg.exe"),
		}
	default:
		return "", fmt.Errorf("unsupported tool: %s", toolName)
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}

		for _, match := range matches {
			if info, err := os.Stat(match); err == nil && !info.IsDir() {
				return match, nil
			}
		}
	}

	return "", fmt.Errorf("%s not found", toolName)
}
