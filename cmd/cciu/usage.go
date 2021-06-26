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
