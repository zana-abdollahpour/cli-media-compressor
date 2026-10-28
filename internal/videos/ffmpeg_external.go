//go:build !embedffmpeg || !windows || !amd64

package videos

import (
	"fmt"
	"os/exec"
)

func ffmpegExecutable() (string, error) {
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		return "", fmt.Errorf("ffmpeg is required to compress videos (install it and add it to PATH): %w", err)
	}
	return path, nil
}
