package main

import (
	"context"
	"os"

	"git.sr.ht/~shinyzenith/devbox/pkg/imageutil"
	"github.com/alecthomas/kong"
	"github.com/containerd/console"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/cmd/ctr/commands"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/sirupsen/logrus"
)

type CLI struct {
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
	context        context.Context
	args           CLI
}

func main() {
	var cli CLI

	// Parsing command line args
	kong_ctx := kong.Parse(&cli, kong.Name("Devbox"), kong.Description("Easy, simple, ephemeral dev containers"), kong.UsageOnError(), kong.ConfigureHelp(kong.HelpOptions{Compact: true, Summary: true}))
	switch kong_ctx.Command() {
	case "create <image> <tag> <name>":
		CFG.image_name = cli.Create.Image
		CFG.image_tag = cli.Create.Tag
		CFG.container_name = cli.Create.Name
		CFG.args = cli
		CFG.context = namespaces.WithNamespace(context.Background(), "Devbox")
	}

	// Set log level
	if cli.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Creating namespace
	logrus.Debug("Creating Devbox namespace with background context")
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		logrus.Errorf("%s", err)
		os.Exit(1)
	}
	defer client.Close()

	// Fetching container image
	image, err := imageutil.GetImage(CFG.context, client, CFG.image_name, CFG.image_tag)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// Creating the container
	logrus.Debug("Creating the container")
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Error(err)
	}
	container, err := client.NewContainer(
		CFG.context,
		CFG.container_name,
		containerd.WithNewSnapshot(CFG.container_name, image),
		containerd.WithNewSpec(oci.WithImageConfig(image), oci.WithHostname(hostname)),
	)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	defer container.Delete(CFG.context, containerd.WithSnapshotCleanup)

	con := console.Current()
	defer con.Reset()
	if err := con.SetRaw(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	if oldTask, err := container.Task(CFG.context, nil); err == nil {
		if _, err := oldTask.Delete(CFG.context); err != nil {
			logrus.Error(err)
		}
	}
	task, err := container.NewTask(CFG.context, cio.NewCreator(cio.WithStreams(con, con, nil), cio.WithTerminal))
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	defer task.Delete(CFG.context)

	exitStatus, err := task.Wait(CFG.context)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	if err := task.Start(CFG.context); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	sigc := commands.ForwardAllSignals(CFG.context, task)
	defer commands.StopCatch(sigc)
	status := <-exitStatus
	code, _, err := status.Result()
	if err != nil {
		logrus.Error(err)
	}
	logrus.Infof("Exit code %d", code)
}
