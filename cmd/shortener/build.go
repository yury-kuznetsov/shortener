package main

import "fmt"

var buildVersion string
var buildDate string
var buildCommit string

func printBuildData() {
	fmt.Printf("Build version: %s\n", getValue(buildVersion))
	fmt.Printf("Build date: %s\n", getValue(buildDate))
	fmt.Printf("Build commit: %s\n", getValue(buildCommit))
}

func getValue(param string) string {
	if param == "" {
		return "N/A"
	}

	return param
}
