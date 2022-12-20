package main

import "go-graph/cmd"

func main() {
	cmd.Init(&cmd.AppInfo{
		Commit:  "",
		Build:   1,
		Name:    "payment service",
		Version: "v1",
		Usage:   "handle payment service",
	})
}
