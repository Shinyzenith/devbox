package main

import (
	. "git.sr.ht/~shinyzenith/devbox/pkg/containerutil"
	. "git.sr.ht/~shinyzenith/devbox/pkg/imageutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newCreateCommand() *cobra.Command {
	var create_command = &cobra.Command{
		Use:   "create IMAGE_NAME IMAGE_TAG CONTAINER_NAME",
		Args:  cobra.MinimumNArgs(3),
		Short: "Create a dev container",
		RunE:  createAction,
	}
	create_command.Flags().SetInterspersed(false)
	return create_command
}

func createAction(cmd *cobra.Command, args []string) error {
	logrus.Info("Creating containerd client!")
	client, err := CreateClient(cmd)
	if err != nil {
		return err
	}
	defer client.Close()

	ctx, err := GetContainerdContext(cmd)
	if err != nil {
		return err
	}

	logrus.Infof("Attempting to fetch %s:%s image", args[0], args[1])
	image, err := GetImage(ctx, client, args[0], args[1])
	if err != nil {
		return err
	}

	container, err := CreateContainer(cmd, client, ctx, args[2], image)
	if err != nil {
		return err
	}
	_ = container

	return nil
}
