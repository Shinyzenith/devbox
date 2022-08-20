package main

import (
	"context"
	"fmt"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func main() {
	ParseFlags()
	ctx := namespaces.WithNamespace(context.Background(), "Devbox")

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	image, err := GetImage(ctx, client, "alpine", "latest")
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
