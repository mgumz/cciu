package main

import "fmt"

var (
	versionString = "0.6.0-dev-build"
	gitHash       = ""
	buildDate     = ""
)

func printVersion(asJSON bool) {
	if asJSON {
		printJSONVersion()
	} else {
		printTextVersion()
	}
}

func printJSONVersion() {
	fmt.Println("{")
	if gitHash != "" {
		fmt.Printf("  %q: %q,\n", "gitHash", gitHash)
	}
	if buildDate != "" {
		fmt.Printf("  %q: %q,\n", "buildDate", buildDate)
	}
	fmt.Printf("  %q: %q\n", "version", versionString)
	fmt.Println("}")
}

func printTextVersion() {
	fmt.Println("cciu:\t" + versionString)
	if gitHash != "" {
		fmt.Println("git:\t" + gitHash)
	}
	if buildDate != "" {
		fmt.Println("build:\t" + buildDate)
	}
}
