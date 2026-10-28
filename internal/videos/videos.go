package videos

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var supportedExtensions = map[string]bool{
	".avi":  true,
	".m4v":  true,
	".mkv":  true,
	".mov":  true,
	".mp4":  true,
	".webm": true,
}

// Supports reports whether the filename extension is handled as a video.
func Supports(ext string) bool {
	return supportedExtensions[strings.ToLower(ext)]
}

// File recompresses a video with FFmpeg. It keeps the source if encoding does
// not reduce its size. FFmpeg must be available on PATH.
func File(src, dst string) (int64, int64, error) {
	ffmpeg, err := exec.LookPath("ffmpeg")
	if err != nil {
		return 0, 0, fmt.Errorf("ffmpeg is required to compress videos (install it and add it to PATH): %w", err)
	}
	inInfo, err := os.Stat(src)
	if err != nil {
		return 0, 0, err
	}

	ext := strings.ToLower(filepath.Ext(dst))
	tmp := strings.TrimSuffix(dst, filepath.Ext(dst)) + ".tmp" + filepath.Ext(dst)
	defer os.Remove(tmp)
	args := []string{"-hide_banner", "-loglevel", "error", "-y", "-i", src, "-map", "0:v:0", "-map", "0:a?"}
	if ext == ".webm" {
		args = append(args, "-c:v", "libvpx-vp9", "-crf", "32", "-b:v", "0", "-c:a", "libopus")
	} else {
		args = append(args, "-c:v", "libx264", "-crf", "28", "-preset", "medium", "-c:a", "aac", "-b:a", "128k")
	}
	if ext == ".mp4" || ext == ".m4v" || ext == ".mov" {
		args = append(args, "-movflags", "+faststart")
	}
	args = append(args, tmp)
	if output, err := exec.Command(ffmpeg, args...).CombinedOutput(); err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			return 0, 0, fmt.Errorf("ffmpeg: %w", err)
		}
		return 0, 0, fmt.Errorf("ffmpeg: %s", message)
	}

	outInfo, err := os.Stat(tmp)
	if err != nil {
		return 0, 0, err
	}
	if outInfo.Size() >= inInfo.Size() {
		if err := copyOriginal(src, dst); err != nil {
			return 0, 0, err
		}
		outInfo, err = os.Stat(dst)
		if err != nil {
			return 0, 0, err
		}
	} else if err := os.Rename(tmp, dst); err != nil {
		return 0, 0, err
	}
	return inInfo.Size(), outInfo.Size(), nil
}

func copyOriginal(src, dst string) error {
	input, err := os.Open(src)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
