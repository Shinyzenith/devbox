# Devbox

<p align=center>
	<a href="https://github.com/shinyzenith/devbox/actions"><img src="https://github.com/shinyzenith/devbox/actions/workflows/ubuntu.yaml/badge.svg"></a>
    <a href="https://goreportcard.com/report/github.com/shinyzenith/devbox"><img src="https://goreportcard.com/badge/github.com/shinyzenith/devbox"></a>
</p>

## Building:

```bash
$ make
# make install
```

## How do I use this?
```bash
$ # Start the containerd daemon.
# containerd
# devbox --debug run --network --rm alpine edge
```

## Requirements:

1. [containerd](https://github.com/containerd/containerd)
1. [go](https://github.com/golang/go)
1. [make](https://git.savannah.gnu.org/cgit/make.git)
1. [musl-gcc](https://www.musl-libc.org/) (Optional: Used to compile statically linked binaries.)
1. [zig](https://github.com/ziglang/zig) (Optional: Used to compile statically linked binaries.)
