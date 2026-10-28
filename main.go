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
	flag.IntVar(&cfg.quality, "quality", 40, "JPEG quality (1-100), higher = better quality")
	flag.BoolVar(&cfg.verbose, "v", false, "verbose output")
	flag.Parse()
	cfg.qualityExplicit = false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "quality" {
			cfg.qualityExplicit = true
		}
	})
	return cfg
}
