package main

import (
	"github.com/ghchinoy/cectl/cmd"
)

func main() {
	cectl := cmd.RootCmd
	cectl.GenBashCompletionFile("cectl.sh")
}
