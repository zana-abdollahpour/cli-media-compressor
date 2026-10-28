package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	cfg := parseFlags()
	if err := run(cfg, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseFlags() config {
	cfg := config{}
	flag.StringVar(&cfg.outputDir, "out", "compressed", "output directory")
	flag.IntVar(&cfg.imageQuality, "image-quality", defaultQuality, "image quality (1-100), higher = better quality")
	flag.IntVar(&cfg.videoQuality, "video-quality", defaultQuality, "video quality (1-100), higher = better quality")
	flag.IntVar(&cfg.quality, "quality", defaultQuality, "set both image and video quality (1-100)")
	flag.BoolVar(&cfg.verbose, "v", false, "verbose output")
	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "image-quality":
			cfg.imageQualitySet = true
		case "video-quality":
			cfg.videoQualitySet = true
		case "quality":
			cfg.qualityExplicit = true
		}
	})
	if cfg.qualityExplicit {
		if !cfg.imageQualitySet {
			cfg.imageQuality = cfg.quality
		}
		if !cfg.videoQualitySet {
			cfg.videoQuality = cfg.quality
		}
	}
	return cfg
}
