package main

import (
	"flag"
	"fmt"
	"os"
)

func printUsage() {

	fmt.Fprintln(os.Stderr, `cciu - check container images for updates

Usage: cciu [flags] <image:1> <image:2> ...
Flags:`)

	flag.PrintDefaults()
}

func printUnsupportedMinMajorLevel(level string) int {

	fmt.Fprintf(os.Stderr, "Ignoring unknown version Level: %s\n", level)
	return 13
}
