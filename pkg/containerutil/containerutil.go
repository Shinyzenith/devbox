package containerutil

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/spf13/cobra"
)

func CreateClient(cmd *cobra.Command) (*containerd.Client, error) {
	if socket, err := cmd.Flags().GetString("socket"); err != nil {
		return nil, err
	} else {
		client, err := containerd.New(socket)
		if err != nil {
			return nil, fmt.Errorf("Failed to connect to containerd socket: %s", err.Error())
		}
		return client, nil
	}
}

func GetContainerdContext(cmd *cobra.Command) (context.Context, error) {
	if namespace, err := cmd.Flags().GetString("namespace"); err != nil {
		return nil, err
	} else {
		return namespaces.WithNamespace(context.Background(), namespace), nil
	}
}

func CreateContainer(cmd *cobra.Command, client *containerd.Client, ctx context.Context, name string, image containerd.Image) (containerd.Container, error) {
	hostname, err := cmd.Flags().GetString("hostname")
	if err != nil {
		return nil, err
	}
	container, err := client.NewContainer(
		ctx,
		name,
		containerd.WithNewSnapshot(name, image),
		containerd.WithNewSpec(oci.WithImageConfig(image), oci.WithHostname(hostname)),
	)
	if err != nil {
		return nil, err
	} else {
		return container, nil
	}
}
