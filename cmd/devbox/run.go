package main

import (
	"fmt"

	. "git.sr.ht/~shinyzenith/devbox/pkg/containerutil"
	. "git.sr.ht/~shinyzenith/devbox/pkg/imageutil"
	"github.com/containerd/console"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cmd/ctr/commands/tasks"
	"github.com/containerd/containerd/contrib/seccomp"
	"github.com/containerd/containerd/oci"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func newRunCommand() *cobra.Command {
	var run_command = &cobra.Command{
		Use:   "run IMAGE_NAME IMAGE_TAG",
		Args:  cobra.MinimumNArgs(2),
		Short: "Run a dev container",
		RunE:  runAction,
	}
	run_command.Flags().SetInterspersed(false)
	setRunFlags(run_command)
	return run_command
}

func setRunFlags(cmd *cobra.Command) {
	//TODO: make use of this lol
	cmd.Flags().BoolP("tty", "t", false, "Attach to container iostream")
	cmd.Flags().Bool("rm", false, "Remove container on successful execution")
}

func runAction(cmd *cobra.Command, args []string) error {
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

	image, err := GetImage(ctx, client, args[0], args[1])
	if err != nil {
		return err
	}

	var (
		opts  []oci.SpecOpts
		cOpts []containerd.NewContainerOpts
	)

	ref_id, err := GenID()
	if err != nil {
		return err
	}

	// TODO: https://github.com/containerd/containerd/blob/main/contrib/nvidia/nvidia.go maybe support gpu passthrough for game emulation via wayland socket exposure??
	// TODO: NETWORK! ( go-cni or libnetwork? libnetwork comes with containerd btw )
	// TODO: Get volumes to work!
	// TODO: Enabling exposing ports!
	// TODO: Support limiting CPU, RAM, and crgoups
	// TODO: Runtime editing of container resources
	// TODO: Inherit cgroups
	// TODO: pid limits
	hostname, err := cmd.Flags().GetString("hostname")
	if err != nil {
		return err
	}

	// Creating container opts
	opts = append(opts, oci.WithDefaultSpec(), oci.WithDefaultUnixDevices, oci.WithDefaultPathEnv, oci.WithImageConfig(image), seccomp.WithDefaultProfile(), oci.WithTTY, oci.WithHostname(hostname))
	cOpts = append(cOpts, containerd.WithImage(image), containerd.WithNewSnapshot(ref_id, image), containerd.WithImageStopSignal(image, "SIGTERM"))
	cOpts = append(cOpts, containerd.WithNewSpec(opts...))

	// Creating container
	container, err := CreateContainer(client, ctx, ref_id, image, cOpts)
	if err != nil {
		return err
	}

	// Remove on execution
	if rm, err := cmd.Flags().GetBool("rm"); err != nil {
		return err
	} else {
		if rm {
			defer container.Delete(ctx, containerd.WithSnapshotCleanup)
		}
	}

	// Setting up console
	console := console.Current()
	defer console.Reset()

	if err := console.SetRaw(); err != nil {
		return err
	}

	task, err := tasks.NewTask(ctx, client, container, "", console, false, "", nil)
	if err != nil {
		return err
	}
	defer task.Delete(ctx)

	statusC, err := task.Wait(ctx)
	if err != nil {
		return err
	}

	if err := task.Start(ctx); err != nil {
		return err
	}

	if err := tasks.HandleConsoleResize(ctx, task, console); err != nil {
		logrus.Errorf("Failed to handle resize console: %s", err)
	}

	status := <-statusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}

	if code != 0 {
		return fmt.Errorf("Container task exited with non-zero exit code: %d", code)
	}
	return nil
}
