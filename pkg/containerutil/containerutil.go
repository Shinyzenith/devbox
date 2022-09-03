package containerutil

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
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

func CreateContainer(client *containerd.Client, ctx context.Context, ref_id string, image containerd.Image, cOpts []containerd.NewContainerOpts) (containerd.Container, error) {
	container, err := client.NewContainer(ctx, ref_id, cOpts...)
	if err != nil {
		return nil, err
	} else {
		return container, nil
	}
}

func GenID() (string, error) {
	h := sha256.New()
	if err := binary.Write(h, binary.LittleEndian, time.Now().UnixNano()); err != nil {
		return "", err
	} else {
		return hex.EncodeToString(h.Sum(nil)), nil
	}
}
