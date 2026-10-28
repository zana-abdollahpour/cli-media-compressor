package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"imgcompress/internal/compress"
)

type config struct {
	outputDir       string
	quality         int
	qualityExplicit bool
	verbose         bool
}

func run(cfg config, in io.Reader, out, errOut io.Writer) error {
	quality, err := resolveQuality(cfg, in, out, errOut)
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting cwd: %w", err)
	}
	outAbs := filepath.Join(cwd, cfg.outputDir)
	if err := os.MkdirAll(outAbs, 0755); err != nil {
		return fmt.Errorf("error creating output dir: %w", err)
	}

	fmt.Fprintf(out, "Scanning: %s\n", cwd)
	fmt.Fprintf(out, "Output:   %s\n\n", outAbs)

	var totalIn, totalOut int64
	var count, skipped int
	err = filepath.Walk(cwd, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			if path == outAbs || strings.HasPrefix(path, outAbs+string(os.PathSeparator)) {
				return filepath.SkipDir
			}
			return nil
		}
		if !compress.Supports(filepath.Ext(path)) {
			return nil
		}
		rel, err := filepath.Rel(cwd, path)
		if err != nil {
			return err
		}
		dst := filepath.Join(outAbs, rel)
		if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
			return err
		}
		inSize, outSize, err := compress.File(path, dst, quality)
		if err != nil {
			fmt.Fprintf(errOut, "  ✗ %s: %v\n", rel, err)
			skipped++
			return nil
		}
		totalIn += inSize
		totalOut += outSize
		count++
		if cfg.verbose {
			saved := inSize - outSize
			pct := 0.0
			if inSize > 0 {
				pct = float64(saved) / float64(inSize) * 100
			}
			fmt.Fprintf(out, "  ✓ %-50s %8d → %8d  (%+.1f%%)\n", rel, inSize, outSize, -pct)
		} else {
			fmt.Fprintf(out, "  ✓ %s\n", rel)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("walk error: %w", err)
	}

	fmt.Fprintf(out, "\nDone. %d files compressed", count)
	if skipped > 0 {
		fmt.Fprintf(out, ", %d skipped", skipped)
	}
	fmt.Fprintln(out)
	if totalIn > 0 {
		saved := totalIn - totalOut
		pct := float64(saved) / float64(totalIn) * 100
		fmt.Fprintf(out, "Total: %d → %d bytes (saved %d bytes, %.1f%%)\n", totalIn, totalOut, saved, pct)
	}
	return nil
}
