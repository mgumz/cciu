package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	usage = `cciu - check container images for updates

Usage: cciu [flags] <image:1> <image:2> ...
Flags:`
)

func printUsage() {

	fmt.Fprintln(os.Stderr, usage)
	flag.PrintDefaults()
}

func printUnsupportedMinMajorLevel(level string) int {

	const unlucky = 13

	fmt.Fprintf(os.Stderr, "Ignoring unknown version Level: %s\n", level)
	return unlucky
}
