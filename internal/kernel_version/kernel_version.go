package kversion

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const (
	STABLE_VERSION_REGEXP = `([0-9]+)\.([0-9]+)\.?([0-9]*)`
)

var (
	MIN_KERNEL_VERSION = KernelVersion{Major: 6, Minor: 16, Patch: -1}
)

type KernelVersion struct {
	Major int8  // > 0
	Minor int8  // -1 field is empty
	Patch int16 // -1 field is empty
}

func (v KernelVersion) String() string {
	if v.Patch < 0 {
		return fmt.Sprintf("%d.%d", v.Major, v.Minor)
	} else {
		return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
}

func Compare(v1, v2 KernelVersion) int {
	if v1.Major > v2.Major {
		return 1
	}
	if v1.Major < v2.Major {
		return -1
	}
	// majors are equal
	if v1.Minor > v2.Minor {
		return 1
	}
	if v1.Minor < v2.Minor {
		return -1
	}
	// majors and minors are equal
	if v1.Patch < 0 && v2.Patch < 0 {
		return 0
	}
	if v1.Patch >= 0 && v2.Patch < 0 {
		return 1
	}
	if v1.Patch < 0 && v2.Patch >= 0 {
		return -1
	}
	if v1.Patch > v2.Patch {
		return 1
	}
	if v1.Patch < v2.Patch {
		return -1
	}
	return 0
}

func Parse(version string) (KernelVersion, error) {
	stableVersionRegExp := regexp.MustCompile(STABLE_VERSION_REGEXP)
	versionSlice := stableVersionRegExp.FindStringSubmatch(version)

	major, err := strconv.Atoi(versionSlice[1])
	if err != nil {
		fmt.Println("Error converting string to integer:", err)
		return KernelVersion{}, err
	}
	minor, err := strconv.Atoi(versionSlice[2])
	if err != nil {
		minor = -1
	}
	patch, err := strconv.Atoi(versionSlice[3])
	if err != nil {
		patch = -1
	}
	parsedVersion := KernelVersion{Major: int8(major), Minor: int8(minor), Patch: int16(patch)}

	return parsedVersion, nil
}

func LatestSystemInstalledVersion() (KernelVersion, error) {
	res := KernelVersion{}
	cmd := exec.Command("dpkg-query",
		"-f=${db:Status-Status} ${Version} ${source:Upstream-Version} ${binary:Package}\n",
		"-W", "linux-image-*")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing dpkg-query: %v\n", err)
		return res, err
	}

	reader := bytes.NewReader(output)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " \t\n")
		pkgData := strings.Split(line, " ")
		if len(pkgData) < 4 {
			continue
		}
		if pkgData[0] != "installed" {
			continue
		}
		currVersion, err := Parse(pkgData[1])
		if err != nil {
			fmt.Printf("Error parsing version %s: %v\n", pkgData[1], err)
		}
		if Compare(res, currVersion) < 0 {
			res = currVersion
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning buffer: %v", err)
	}

	return res, nil
}
