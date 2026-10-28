package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const defaultQuality = 40

func resolveQuality(cfg config, in io.Reader, out, errOut io.Writer) (int, error) {
	if cfg.qualityExplicit {
		if cfg.quality < 1 || cfg.quality > 100 {
			return 0, errors.New("Error: quality must be between 1 and 100.")
		}
		return cfg.quality, nil
	}

	fmt.Fprint(out, "Enter JPEG quality (1-100, default 40): ")
	var input string
	if _, err := fmt.Fscanln(in, &input); err != nil {
		if errors.Is(err, io.EOF) {
			return defaultQuality, nil
		}
		if err.Error() != "unexpected newline" {
			return 0, err
		}
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultQuality, nil
	}
	q, err := strconv.Atoi(input)
	if err != nil || q < 1 || q > 100 {
		return 0, errors.New("Error: quality must be a number between 1 and 100.")
	}
	return q, nil
}
