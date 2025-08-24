// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright 2025 Alex Syrnikov <alex19srv@gmail.com>
package cmd

import (
	"fmt"
	kversion "mainline-kernel-tool/internal/kernel_version"
	mrepo "mainline-kernel-tool/internal/ubuntu-mainline-repo"

	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [latest_installed_version]",
	Short: "check new ubuntu mainline kernel versions",
	Long: `Check for new ubuntu mainline kernel, newer then provided
in command line. If version in command line is omited,
the latest installed version is used (via dpkg-query).
If no new kernel is found, will be used default minimal kernel version.
Example:
	mainline-kernel check 6.16.1`,
	Run: checkLatestVersion,
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func checkLatestVersion(cmd *cobra.Command, args []string) {
	latestVersion, err := mrepo.LatestVersion()
	if err != nil {
		fmt.Printf("failed get LatestVersion(): %v\n", err)
		return
	}

	var cmdLineVersion string
	if len(args) > 0 {
		cmdLineVersion = args[0]
	}

	currentVersion, err := installedVersion(cmdLineVersion)
	if err != nil {
		fmt.Printf("error in installedVersion(): %v\n", err)
		return
	}
	if kversion.Compare(latestVersion, currentVersion) > 0 {
		fmt.Printf("latest version %v\n", latestVersion)
	}
}
