/*
* SPDX-License-Identifier: GPL-3.0-only
*
* main.go
*
* Created by:	Aakash Sen Sharma
* Copyright:	(C) 2022, Aakash Sen Sharma & Contributors
 */

package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:              "devbox",
		Short:            "Instant, easy, ephemeral development containers.",
		Long:             "Instant, easy, ephemeral development containers.",
		Version:          "0.1.0",
		TraverseChildren: true,
		SilenceUsage:     true,
	}
)

func main() {
	// Setting flags!
	rootCmd.PersistentFlags().StringP("namespace", "n", "devbox", "Set the containerd namespace for usage")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode")

	//TODO: Maybe fallback to /run/user/$(user_id)/containerd.sock if the default one doesn't exist?
	rootCmd.PersistentFlags().StringP("socket", "s", "/run/containerd/containerd.sock", "Set the containerd socket")

	if hostname, err := os.Hostname(); err != nil {
		logrus.Fatal(err)
	} else {
		rootCmd.PersistentFlags().String("hostname", hostname, "Set the container hostname")
	}

	// Setting pre-run hook
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		debug, err := cmd.Flags().GetBool("debug")
		if err != nil {
			logrus.Fatal(err)
		}

		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}
	}

	// Adding all the subcommands
	rootCmd.AddCommand(
		newRunCommand(),
	)

	// Parsing flags
	rootCmd.Execute()
}
