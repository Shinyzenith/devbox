package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Create struct {
		Image string `arg:"" help:"Set the container image" type:"string" name:"image"`
		Tag   string `arg:"" help:"Set the image tag." type:"string" name:"tag"`

		Name string `arg:"" name:"name" help:"Set the container name." type:"path"`
	} `cmd:"" help:"Create a dev container."`
}

func ParseFlags() {
	kong_ctx := kong.Parse(&CLI, kong.Name("Devbox"), kong.Description("Easy, simple, ephemeral dev containers"), kong.UsageOnError(), kong.ConfigureHelp(kong.HelpOptions{Compact: true, Summary: true}))
	switch kong_ctx.Command() {
	case "create <image> <tag> <name>":
		fmt.Println(CLI.Create.Image)
	}
}
