package main

import (
	"context"
	"os"

	"github.com/alecthomas/kong"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/sirupsen/logrus"
)

var CLI struct {
	Debug  bool `help:"Enable debug mode"`
	Create struct {
		Image string `arg:"" help:"Set the container image" type:"string" name:"image"`
		Tag   string `arg:"" help:"Set the image tag." type:"string" name:"tag"`

		Name string `arg:"" name:"name" help:"Set the container name." type:"string"`
	} `cmd:"" help:"Create a dev container."`
}

var CFG struct {
	image_name     string
	image_tag      string
	container_name string
}

func main() {
	// Parsing command line args
	kong_ctx := kong.Parse(&CLI, kong.Name("Devbox"), kong.Description("Easy, simple, ephemeral dev containers"), kong.UsageOnError(), kong.ConfigureHelp(kong.HelpOptions{Compact: true, Summary: true}))
	switch kong_ctx.Command() {
	case "create <image> <tag> <name>":
		CFG.image_name = CLI.Create.Image
		CFG.image_tag = CLI.Create.Tag
		CFG.container_name = CLI.Create.Name
	}

	// Set log level
	if CLI.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Infof("The following config was detected: %s", CFG)

	// Creating namespace
	logrus.Debug("Creating Devbox namespace with background context")
	ctx := namespaces.WithNamespace(context.Background(), "Devbox")
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		logrus.Errorf("%s", err)
		os.Exit(1)
	}
	defer client.Close()

	// Fetching container image
	image, err := GetImage(ctx, client, CFG.image_name, CFG.image_tag)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// Creating the container
	logrus.Debug("Creating the container")
	container, err := client.NewContainer(
		ctx,
		"replace_me_later",
		// By using different id's we can reuse them across containers
		containerd.WithNewSnapshot("replace_me_later", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)
}

func GetImage(ctx context.Context, client *containerd.Client, image_name string, image_tag string) (containerd.Image, error) {
	logrus.Debugf("Fetching %s:%s image", image_name, image_tag)
	image_url, err := ResolveShortnameUrl(image_name, image_tag)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	image, err := client.GetImage(ctx, image_url)
	if err != nil {
		logrus.Debugf("%s:%s not found, Pulling from source.\n", image_name, image_tag)
		image, err := client.Pull(ctx, image_url, containerd.WithPullUnpack)
		if err != nil {
			return nil, err
		}
		return image, nil
	}
	logrus.Debugf("Found %s:%s image, not pulling.\n", image_name, image_tag)
	return image, nil
}
