//go:build !ndi

package main

import "fmt"

// printUsage prints the client usage for builds without NDI support. It
// intentionally contains no NDI flags, attribution, or links.
func printUsage() {
	fmt.Print(`Usage: client [options] [ws://host:port] [ClientName]

This build does not include NDI support.
`)
}
