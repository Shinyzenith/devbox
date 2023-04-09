# Devbox

<p align=center>
    <a href="https://builds.sr.ht/~shinyzenith/devbox/commits/master/arch.yml"><img src="https://builds.sr.ht/~shinyzenith/devbox/commits/master/arch.yml.svg"</a>
	<a href="https://github.com/shinyzenith/devbox/actions"><img src="https://github.com/shinyzenith/devbox/actions/workflows/ubuntu.yaml/badge.svg"></a>
    <a href="https://goreportcard.com/report/git.sr.ht/~shinyzenith/devbox"><img src="https://goreportcard.com/badge/git.sr.ht/~shinyzenith/devbox"></a>
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
