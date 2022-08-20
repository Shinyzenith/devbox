package main

import (
	"context"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

var CLI struct {
	Create struct {
		Image string `arg:"" help:"Set the container image" type:"string" name:"image"`
		Tag   string `arg:"" help:"Set the image tag." type:"string" name:"tag"`

		Name string `arg:"" name:"name" help:"Set the container name." type:"path"`
	} `cmd:"" help:"Create a dev container."`
}

var CFG struct {
	image_name     string
	image_tag      string
	container_name string
}

func main() {
	kong_ctx := kong.Parse(&CLI, kong.Name("Devbox"), kong.Description("Easy, simple, ephemeral dev containers"), kong.UsageOnError(), kong.ConfigureHelp(kong.HelpOptions{Compact: true, Summary: true}))
	switch kong_ctx.Command() {
	case "create <image> <tag> <name>":
		CFG.image_name = CLI.Create.Image
		CFG.image_tag = CLI.Create.Tag
		CFG.container_name = CLI.Create.Name
	}

	ctx := namespaces.WithNamespace(context.Background(), "Devbox")
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	image, err := GetImage(ctx, client, CFG.image_name, CFG.image_tag)
	if err != nil {
		panic(err)
	}

	container, err := client.NewContainer(
		ctx,
		"replace_me_later",
		// By using different id's we can reuse them across containers
		containerd.WithNewSnapshot("replace_me_later", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		panic(err)
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)
}

func GetImage(ctx context.Context, client *containerd.Client, image_name string, image_tag string) (containerd.Image, error) {
	image_url, err := ResolveShortnameUrl(image_name, image_tag)
	if err != nil {
		panic(err)
	}

	image, err := client.GetImage(ctx, image_url)
	if err != nil {
		fmt.Printf("%s:%s not found, Pulling from source.\n", image_name, image_tag)
		image, err := client.Pull(ctx, image_url, containerd.WithPullUnpack)
		if err != nil {
			return nil, err
		}
		return image, nil
	}
	fmt.Printf("Found %s:%s image, not pulling.\n", image_name, image_tag)
	return image, nil
}
