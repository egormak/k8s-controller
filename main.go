package main

import "k8s-controller/cmd"

// version will be set during build time via ldflags
var version = "dev"

func main() {
	cmd.SetVersion(version)
	cmd.Execute()
}
