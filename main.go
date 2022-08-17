package main

import (
	"context"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func main() {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "Devbox")
	image, err := client.Pull(ctx, "docker.io/library/alpine:latest", containerd.WithPullUnpack)
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
