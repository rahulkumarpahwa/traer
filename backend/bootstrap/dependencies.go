package bootstrap

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type DependencyManager struct {
	Required map[string]string
}

func NewDependencyManager() *DependencyManager {
	return &DependencyManager{
		Required: make(map[string]string),
	}
}

func (d *DependencyManager) Ensure(tools []string) error {
	var missing []string

	for _, tool := range tools {
		// 1) try PATH
		path, err := exec.LookPath(tool)
		if err == nil {
			d.Required[tool] = path
			log.Printf("found %s at %s (PATH)", tool, path)
			continue
		}

		// 2) try WinGet / LOCALAPPDATA locations
		if p, rerr := resolveToolPath(tool); rerr == nil {
			d.Required[tool] = p
			log.Printf("found %s at %s (resolved)", tool, p)
			continue
		}

		// 3) attempt install (fallback)
		log.Printf("tool %s not found, attempting install...", tool)
		if installErr := installTool(tool); installErr != nil {
			// installation failed; mark missing
			missing = append(missing, tool)
			continue
		}

		// retry after install: first PATH then resolve locations
		path, err = exec.LookPath(tool)
		if err == nil {
			d.Required[tool] = path
			log.Printf("found %s at %s (after install PATH)", tool, path)
			continue
		}
		if p, rerr := resolveToolPath(tool); rerr == nil {
			d.Required[tool] = p
			log.Printf("found %s at %s (after install resolved)", tool, p)
			continue
		}

		missing = append(missing, tool)
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing dependencies: %v", missing)
	}

	return nil
}

// resolveToolPath searches common install locations (WinGet) for the given tool.
func resolveToolPath(toolName string) (string, error) {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return "", fmt.Errorf("LOCALAPPDATA not set")
	}

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
		for _, m := range matches {
			if info, err := os.Stat(m); err == nil && !info.IsDir() {
				return m, nil
			}
		}
	}

	return "", fmt.Errorf("%s not found", toolName)
}
