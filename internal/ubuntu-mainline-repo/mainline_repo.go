package mainline_repo

import (
	"bufio"
	"fmt"
	"io"
	"mainline-kernel-tool/internal/http_client"
	kversion "mainline-kernel-tool/internal/kernel_version"
	"os"
	"regexp"
	"slices"
	"strconv"
)

const (
	UBUNTU_MAINLINE_URL = "https://kernel.ubuntu.com/mainline/"
)

type netClient interface {
	Get(url string) ([]byte, error)
	GetAndSave(url string, writer io.Writer) error
}

type MainlineKernelRepo struct {
	client netClient
}

func New() *MainlineKernelRepo {
	return &MainlineKernelRepo{
		client: http_client.New(),
	}
}
func (repo *MainlineKernelRepo) NewVersions(filterVersion *kversion.KernelVersion) ([]kversion.KernelVersion, error) {
	if filterVersion == nil {
		filterVersion = &kversion.MIN_KERNEL_VERSION
	}
	body, err := repo.client.Get(UBUNTU_MAINLINE_URL)
	if err != nil {
		fmt.Println("error in getURL():", err)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("body is empty")
	}
	// need match this
	// <a href="v6.16/">
	// <a href="v6.15.9/">
	// <a  href="v6.15.9/">
	// but do not match
	// <a href="v6.16-rc7/">
	stableVersionRegExp := regexp.MustCompile(`<a +href="v([0-9]+)\.([0-9]+)\.?([0-9]*)/">`)

	versions := stableVersionRegExp.FindAllSubmatch(body, -1)
	// versions will contain
	// [ ["<a href=\"v6.15.7/\">" "6" "15" "7"], ["<a href=\"v6.16/\">" "6" "16" ""]]
	res := make([]kversion.KernelVersion, 0)
	for _, v := range versions {
		major, err := strconv.Atoi(string(v[1]))
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			continue
		}
		minor, err := strconv.Atoi(string(v[2]))
		if err != nil {
			minor = -1
		}
		patch, err := strconv.Atoi(string(v[3]))
		if err != nil {
			patch = -1
		}
		currentVersion := kversion.KernelVersion{Major: int8(major), Minor: int8(minor), Patch: int16(patch)}
		if filterVersion != nil {
			if kversion.Compare(currentVersion, *filterVersion) < 0 {
				continue
			}
		}
		res = append(res, currentVersion)
	}
	slices.SortFunc(res, kversion.Compare)

	return res, nil
}

func LatestVersion() (kversion.KernelVersion, error) {
	repo := New()
	return repo.latestVersion()
}
func (repo *MainlineKernelRepo) latestVersion() (kversion.KernelVersion, error) {
	res := kversion.MIN_KERNEL_VERSION

	body, err := repo.client.Get(UBUNTU_MAINLINE_URL)
	if err != nil {
		fmt.Println("error in getURL():", err)
	}
	if len(body) == 0 {
		return res, fmt.Errorf("body is empty")
	}
	// need match this
	// <a href="v6.16/">
	// <a href="v6.15.9/">
	// <a  href="v6.15.9/">
	// but do not match
	// <a href="v6.16-rc7/">
	stableVersionRegExp := regexp.MustCompile(`<a +href="v([0-9]+)\.([0-9]+)\.?([0-9]*)/">`)

	versions := stableVersionRegExp.FindAllSubmatch(body, -1)
	// versions will contain
	// [ ["<a href=\"v6.15.7/\">" "6" "15" "7"], ["<a href=\"v6.16/\">" "6" "16" ""] ]

	for _, v := range versions {
		major, err := strconv.Atoi(string(v[1]))
		if err != nil {
			fmt.Println("Error converting string to integer:", err)
			continue
		}
		minor, err := strconv.Atoi(string(v[2]))
		if err != nil {
			minor = -1
		}
		patch, err := strconv.Atoi(string(v[3]))
		if err != nil {
			patch = -1
		}
		currentVersion := kversion.KernelVersion{Major: int8(major), Minor: int8(minor), Patch: int16(patch)}
		if kversion.Compare(res, currentVersion) < 0 {
			res = currentVersion
		}
	}

	return res, nil
}
func DownloadDebFiles(targetDirectory string, version kversion.KernelVersion) ([]string, error) {
	repo := New()
	return repo.downloadDebFiles(targetDirectory, version)
}

func (repo *MainlineKernelRepo) downloadDebFiles(
	targetDirectory string, version kversion.KernelVersion) ([]string, error) {

	// https://kernel.ubuntu.com/mainline/v6.16.1/amd64/linux-image-unsigned-6.16.1-061601-generic_6.16.1-061601.202508151544_amd64.deb
	versionURL := UBUNTU_MAINLINE_URL + "v" + version.String() + "/"
	html, err := repo.client.Get(versionURL)
	if err != nil {
		fmt.Println("error in get version URL:", err)
	}
	if len(html) == 0 {
		return nil, fmt.Errorf("failed get version URL %s", versionURL)
	}

	// need match this
	// <a href="amd64/linux-image-unsigned-6.16.1-061601-generic_6.16.1-061601.202508151544_amd64.deb">
	// <a href="amd64/linux-modules-6.16.1-061601-generic_6.16.1-061601.202508151544_amd64.deb">
	imageRegExpStr := `<a +href="amd64\/(linux-image-unsigned-` + version.String() + `[0-9a-z-_.]*amd64\.deb)">`
	modulesRegExpStr := `<a +href="amd64\/(linux-modules-` + version.String() + `[0-9a-z-_.]*amd64\.deb)">`

	imageRegExp := regexp.MustCompile(imageRegExpStr)
	modulesRegExp := regexp.MustCompile(modulesRegExpStr)
	images := imageRegExp.FindSubmatch(html)
	modules := modulesRegExp.FindSubmatch(html)

	var imagePath string
	var modulesPath string
	if len(images) != 2 {
		fmt.Printf("len(v) != 2, %d\n", len(images))
		fmt.Printf("failed parse image path\n")
		return nil, nil // FIXME: return error
	}
	// FIXME: check if this is a valid URL (not empty string)
	imageFileName := string(images[1])
	imagePath = versionURL + "amd64/" + string(images[1])
	if len(modules) != 2 {
		fmt.Printf("failed parse modules path\n")
		return nil, nil // FIXME: return error
	}
	// FIXME: check if this is a valid URL (not empty string)
	modulesFileName := string(modules[1])
	modulesPath = versionURL + "amd64/" + string(modules[1])

	res := []string{
		targetDirectory + "/" + imageFileName,
		targetDirectory + "/" + modulesFileName,
	}
	err = repo.downloadAndSave(res[0], imagePath)
	if err != nil {
		fmt.Printf("failed download and save image file %s: %v\n", imagePath, err)
		return nil, err
	}
	err = repo.downloadAndSave(res[1], modulesPath)
	if err != nil {
		fmt.Printf("failed download and save image file %s: %v\n", modulesPath, err)
		return nil, err
	}

	return res, nil
}

func (repo *MainlineKernelRepo) downloadAndSave(debFilePath string, url string) error {
	file, err := os.Create(debFilePath)
	if err != nil {
		fmt.Printf("failed create file %s\n", debFilePath)
		return err
	}
	bufferedFile := bufio.NewWriter(file)
	if err := repo.client.GetAndSave(url, bufferedFile); err != nil {
		fmt.Printf("failed download and save file %s: %v\n", url, err)
		return err
	}
	return nil
}
