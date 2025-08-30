// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright 2025 Alex Syrnikov <alex19srv@gmail.com>
package cmd

import (
	"fmt"

	kversion "github.com/alex19srv/ubuntu-mainline-kernel-tool/internal/kernel_version"
	mrepo "github.com/alex19srv/ubuntu-mainline-kernel-tool/internal/ubuntu-mainline-repo"

	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch-latest [latest_installed_version]",
	Short: "fetch new ubuntu mainline kernel versions",
	Long: `Download latest stable ubuntu mainline kernel, if avalable newer
then provided in command line. If version in command line is omited,
the latest system installed version is used (via dpkg-query).
If no kernel version is found in commandline and system,
will be used default minimal kernel version as currently installed.

!!! You provide already installed version in cmdline and fetch will
download latest avalable !!!

Example:
	mainline-kernel fetch-latest 6.16.1 --dir=/tmp/kernels --force`,
	Run: fetchNewKernelVersion,
}

var (
	targetDirectory string
	forceDownload   bool
)

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("dir", "", "target directory to save debs [default .]")
	fetchCmd.Flags().StringVarP(&targetDirectory, "dir", "", ".", "target directory to save debs")
	fetchCmd.Flags().BoolVarP(&forceDownload, "force", "", false, "force download, even if latest kernel present")
}

func fetchNewKernelVersion(cmd *cobra.Command, args []string) {
	var cmdLineVersion string
	if len(args) > 0 {
		cmdLineVersion = args[0]
	}

	currentVersion, err := installedVersion(cmdLineVersion)
	if err != nil {
		fmt.Printf("error in installedVersion(): %v\n", err)
		return
	}
	latestVersion, err := mrepo.LatestVersion()
	if err != nil {
		fmt.Printf("failed get LatestVersion(): %v\n", err)
		return
	}
	var debFilesPaths []string
	if forceDownload || kversion.Compare(latestVersion, currentVersion) > 0 {
		debFilesPaths, err = mrepo.DownloadDebFiles(targetDirectory, latestVersion)
		if err != nil {
			fmt.Printf("failed download deb files: %v\n", err)
			return
		}
	}
	for _, filePath := range debFilesPaths {
		fmt.Printf("%s\n", filePath)
	}
}
