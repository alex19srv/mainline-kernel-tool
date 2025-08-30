// SPDX-License-Identifier: GPL-3.0-or-later
// Copyright 2025 Alex Syrnikov <alex19srv@gmail.com>
package cmd

import (
	"fmt"

	kversion "github.com/alex19srv/ubuntu-mainline-kernel-tool/internal/kernel_version"
)

func installedVersion(cmdVersion string) (kversion.KernelVersion, error) {
	installedVersion := kversion.MIN_KERNEL_VERSION

	if len(cmdVersion) > 0 {
		v, err := kversion.Parse(cmdVersion)
		if err != nil {
			fmt.Printf("error parsing cmdline version \"%s\" (expected format \"check 6.1.6\"): %v\n",
				cmdVersion, err)
		} else {
			installedVersion = v
			return installedVersion, nil
		}
	}

	v, err := kversion.LatestSystemInstalledVersion()
	if err != nil {
		fmt.Printf("error getting system installed kernel version")
	} else {
		installedVersion = v
	}

	return installedVersion, nil
}
