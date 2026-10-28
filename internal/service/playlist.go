package service

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/rahulkumarpahwa/traer/internal/types"
)

type PlaylistWorker struct {
	ExecutablePaths map[string]string
}

func (pl *PlaylistWorker) RunPlayList(playlist *types.Playlist) error {

	if playlist == nil {
		return fmt.Errorf("playlist must not be nil")
	}

	if playlist.PlaylistUrl == "" {
		return fmt.Errorf("playlist URL must not be empty")
	}

	cmd := pl.buildCommand(playlist)
	if cmd == nil {
		return fmt.Errorf("invalid job type")
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var outputWg sync.WaitGroup
	outputWg.Add(1)
	go func() {
		defer outputWg.Done()
		parseOutput(stdout, playlist)
	}()

	outputWg.Wait()
	waitErr := cmd.Wait()
	if waitErr != nil {
		return waitErr
	}

	// if resolvedPath, err := ResolveOutputPath(.Output); err == nil && resolvedPath != "" {
	// 	job.Output = resolvedPath
	// 	job.TempPath = resolvedPath

	// }
	return nil
}

func (pl *PlaylistWorker) buildCommand(playlist *types.Playlist) *exec.Cmd {
	args := []string{
		"--flat-playlist",
		"--dump-single-json",
		"--skip-download",
		playlist.PlaylistUrl,
	}

	ytDlpPath := pl.ExecutablePaths["yt-dlp"]
	if ytDlpPath == "" {
		return nil
	}

	return exec.Command(ytDlpPath, args...)
}

func parseOutput(pipe io.ReadCloser, playlist *types.Playlist) {
	// scanner := bufio.NewScanner(pipe)
	// buf := make([]byte, 0, 64*1024)
	// scanner.Buffer(buf, 1024*1024)

	json.NewDecoder(pipe).Decode(playlist)
	fmt.Println(playlist.Entries)

	// var lastPathLine string // Capture the last path line to ensure we get the final filename

	// for scanner.Scan() {
	// 	line := scanner.Text()

	// 	// Extract filepath - look for the --print after_move:filepath output
	// 	// yt-dlp outputs just the filepath on its own line when using --print after_move:filepath
	// 	trimmedLine := strings.TrimSpace(line)
	// 	if (strings.HasPrefix(trimmedLine, "C:\\") || strings.HasPrefix(trimmedLine, "/")) &&
	// 		strings.Contains(trimmedLine, "downloads") {
	// 		// This looks like a filepath (Windows or Unix absolute path)
	// 		// Store it, we'll use the last one encountered to ensure we get the final filename
	// 		lastPathLine = trimmedLine
	// 		fmt.Printf("job found path line: %s\n", trimmedLine)
	// 	}
	// }

	// // Set the final output path after processing all output
	// if lastPathLine != "" {
	// 	job.Output = lastPathLine
	// 	fmt.Printf("job Final captured output path: %s\n", lastPathLine)
	// }

	// if err := scanner.Err(); err != nil {
	// 	fmt.Println("Error reading output:", err)
	// }

}

// helpers
// func ResolveOutputPath(outputPath string) (string, error) {
// 	if outputPath == "" {
// 		return "", nil
// 	}

// 	if _, err := os.Stat(outputPath); err == nil {
// 		return outputPath, nil
// 	}

// 	dir := filepath.Dir(outputPath)
// 	base := filepath.Base(outputPath)
// 	ext := filepath.Ext(base)
// 	stem := strings.TrimSuffix(base, ext)
// 	stem = strings.TrimSpace(stem)

// 	entries, err := os.ReadDir(dir)
// 	if err != nil {
// 		return "", err
// 	}

// 	for _, entry := range entries {
// 		if entry.IsDir() {
// 			continue
// 		}

// 		candidateName := entry.Name()
// 		candidateExt := filepath.Ext(candidateName)
// 		if ext != "" && !strings.EqualFold(candidateExt, ext) {
// 			continue
// 		}

// 		candidateStem := strings.TrimSpace(strings.TrimSuffix(candidateName, candidateExt))
// 		if strings.HasPrefix(candidateStem, stem) || strings.HasPrefix(stem, candidateStem) {
// 			return filepath.Join(dir, candidateName), nil
// 		}
// 	}

// 	return "", os.ErrNotExist
// }
