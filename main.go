package main

import (
	"github.com/containerd/containerd"
)

func main() {
	client, err := containerd.New("./test.sock")
	if err != nil {
		panic(err)
	}

	defer client.Close()
}
