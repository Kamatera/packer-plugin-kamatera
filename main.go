package main

import (
	"fmt"
	"os"
	"packer-plugin-kamatera/builder/kamatera"
	"packer-plugin-kamatera/version"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("kamatera-builder", new(kamatera.Builder))
	pps.SetVersion(version.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
