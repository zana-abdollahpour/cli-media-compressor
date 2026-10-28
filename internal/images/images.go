package images

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"
)

var supportedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

// Supports reports whether a filename extension is handled by the compressor.
func Supports(ext string) bool {
	return supportedExtensions[strings.ToLower(ext)]
}

// File rewrites src to dst using the selected format's encoder. If encoding
// does not reduce the file size, it copies the original bytes instead.
func File(src, dst string, quality int) (int64, int64, error) {
	inInfo, err := os.Stat(src)
	if err != nil {
		return 0, 0, err
	}

	in, err := os.Open(src)
	if err != nil {
		return 0, 0, err
	}
	img, format, decodeErr := image.Decode(in)
	closeErr := in.Close()
	if decodeErr != nil {
		return 0, 0, fmt.Errorf("decode: %w", decodeErr)
	}
	if closeErr != nil {
		return 0, 0, closeErr
	}

	out, err := os.Create(dst)
	if err != nil {
		return 0, 0, err
	}

	switch format {
	case "jpeg":
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(out, img)
	case "gif":
		err = gif.Encode(out, img, nil)
	default:
		_ = out.Close()
		_ = os.Remove(dst)
		return 0, 0, fmt.Errorf("unsupported format: %s", format)
	}

	closeErr = out.Close()
	if err != nil {
		_ = os.Remove(dst)
		return 0, 0, fmt.Errorf("encode: %w", err)
	}
	if closeErr != nil {
		_ = os.Remove(dst)
		return 0, 0, closeErr
	}

	outInfo, err := os.Stat(dst)
	if err != nil {
		return 0, 0, err
	}
	if outInfo.Size() >= inInfo.Size() {
		if err := copyFile(src, dst); err != nil {
			return 0, 0, err
		}
		outInfo, err = os.Stat(dst)
		if err != nil {
			return 0, 0, err
		}
	}

	return inInfo.Size(), outInfo.Size(), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(out, in)
	closeErr := out.Close()
	if copyErr != nil {
		return copyErr
	}
	return closeErr
}
