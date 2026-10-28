package main

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const defaultQuality = 40

func resolveQuality(mediaType string, quality int, explicit bool, in io.Reader, out io.Writer) (int, error) {
	label := strings.ToLower(mediaType)
	if explicit {
		if quality < 1 || quality > 100 {
			return 0, fmt.Errorf("Error: %s quality must be between 1 and 100.", label)
		}
		return quality, nil
	}

	fmt.Fprintf(out, "Enter %s quality (1-100, default %d): ", label, defaultQuality)
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
		return 0, fmt.Errorf("Error: %s quality must be a number between 1 and 100.", label)
	}
	return q, nil
}
