//go:build embedffmpeg && windows && amd64

package videos

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed assets/windows_amd64/ffmpeg.exe
var embeddedFFmpeg []byte

func ffmpegExecutable() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("locate FFmpeg cache directory: %w", err)
	}
	installDir := filepath.Join(cacheDir, "imgcompress", "ffmpeg-2026-10-25-git-46f23be8b1")
	if err := os.MkdirAll(installDir, 0700); err != nil {
		return "", fmt.Errorf("create FFmpeg cache directory: %w", err)
	}
	executable := filepath.Join(installDir, "ffmpeg.exe")
	if info, statErr := os.Stat(executable); statErr == nil {
		if info.Size() == int64(len(embeddedFFmpeg)) {
			return executable, nil
		}
		if err := os.Remove(executable); err != nil {
			return "", fmt.Errorf("replace cached FFmpeg executable: %w", err)
		}
	} else if !os.IsNotExist(statErr) {
		return "", fmt.Errorf("check cached FFmpeg executable: %w", statErr)
	}

	tempFile, err := os.CreateTemp(installDir, "ffmpeg-*.exe")
	if err != nil {
		return "", fmt.Errorf("create temporary FFmpeg executable: %w", err)
	}
	tempPath := tempFile.Name()
	defer os.Remove(tempPath)
	if _, err := tempFile.Write(embeddedFFmpeg); err != nil {
		_ = tempFile.Close()
		return "", fmt.Errorf("extract embedded FFmpeg: %w", err)
	}
	if err := tempFile.Close(); err != nil {
		return "", fmt.Errorf("finish extracting embedded FFmpeg: %w", err)
	}
	if err := os.Rename(tempPath, executable); err != nil {
		if info, statErr := os.Stat(executable); statErr == nil && info.Size() == int64(len(embeddedFFmpeg)) {
			return executable, nil
		}
		return "", fmt.Errorf("install embedded FFmpeg: %w", err)
	}
	return executable, nil
}
