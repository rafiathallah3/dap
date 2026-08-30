package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	RepoOwner = "rafiathallah3"
	RepoName  = "dap"
	APIURL    = "https://api.github.com/repos/rafiathallah3/dap/releases/latest"
)

// ReleaseAsset nyimpen info file binary release dari GitHub.
type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// GitHubRelease nyimpen payload response dari GitHub Releases API.
type GitHubRelease struct {
	TagName     string         `json:"tag_name"`
	Name        string         `json:"name"`
	HTMLURL     string         `json:"html_url"`
	PublishedAt string         `json:"published_at"`
	Assets      []ReleaseAsset `json:"assets"`
}

// FetchLatestRelease buat ngambil info rilis terbaru langsung dari API GitHub.
func FetchLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{
		Timeout: 6 * time.Second,
	}

	req, err := http.NewRequest("GET", APIURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "dap-updater")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned status %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decoding release json: %w", err)
	}

	return &release, nil
}

// isNewerVersion returns true if remote version is strictly newer than current version.
func isNewerVersion(remote, current string) bool {
	remote = strings.TrimPrefix(remote, "v")
	current = strings.TrimPrefix(current, "v")
	if current == "dev" {
		return false
	}
	rParts := strings.Split(remote, ".")
	cParts := strings.Split(current, ".")
	maxLen := len(rParts)
	if len(cParts) > maxLen {
		maxLen = len(cParts)
	}
	for i := 0; i < maxLen; i++ {
		var rNum, cNum int
		if i < len(rParts) {
			fmt.Sscanf(rParts[i], "%d", &rNum)
		}
		if i < len(cParts) {
			fmt.Sscanf(cParts[i], "%d", &cNum)
		}
		if rNum > cNum {
			return true
		}
		if rNum < cNum {
			return false
		}
	}
	return false
}

// CheckVersion buat nampilin versi sekarang dan ngecek ada update baru apa engga di GitHub.
func CheckVersion(currentVersion string) {
	fmt.Printf("DAP version %s (%s/%s)\n", currentVersion, runtime.GOOS, runtime.GOARCH)

	fmt.Print("Checking for updates from GitHub... ")
	release, err := FetchLatestRelease()
	if err != nil {
		fmt.Printf("(unable to check for updates: %v)\n", err)
		return
	}

	if isNewerVersion(release.TagName, currentVersion) {
		fmt.Printf("\nA new version of DAP is available: %s -> %s\n", currentVersion, release.TagName)
		fmt.Println("Run 'dap --update' to upgrade.")
	} else {
		fmt.Println("You are on the latest version.")
	}
}

// PerformSelfUpdate buat nge-download dan update binary DAP yang lagi jalan ke versi terbaru.
func PerformSelfUpdate(currentVersion string, force bool) error {
	fmt.Printf("Current DAP version: %s (%s/%s)\n", currentVersion, runtime.GOOS, runtime.GOARCH)
	fmt.Println("Checking for latest release...")

	release, err := FetchLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to fetch release information: %w", err)
	}

	if !force && !isNewerVersion(release.TagName, currentVersion) && currentVersion != "dev" {
		fmt.Printf("DAP is already up to date (%s).\n", release.TagName)
		return nil
	}

	expectedAssetName := getExpectedAssetName()
	if expectedAssetName == "" {
		return fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == expectedAssetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("asset %q not found in release %s", expectedAssetName, release.TagName)
	}

	fmt.Printf("Downloading %s from %s...\n", expectedAssetName, release.TagName)
	archiveData, err := downloadFile(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download release asset: %w", err)
	}

	fmt.Println("Extracting binary...")
	binaryName := "dap"
	if runtime.GOOS == "windows" {
		binaryName = "dap.exe"
	}

	var newBinaryBytes []byte
	if strings.HasSuffix(expectedAssetName, ".zip") {
		newBinaryBytes, err = extractFromZip(archiveData, binaryName)
	} else if strings.HasSuffix(expectedAssetName, ".tar.gz") {
		newBinaryBytes, err = extractFromTarGz(archiveData, binaryName)
	} else {
		return fmt.Errorf("unsupported archive format: %s", expectedAssetName)
	}

	if err != nil {
		return fmt.Errorf("failed to extract binary: %w", err)
	}

	fmt.Println("Replacing executable...")
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to determine executable path: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable symlinks: %w", err)
	}

	if err := replaceExecutable(execPath, newBinaryBytes); err != nil {
		return fmt.Errorf("failed to replace executable: %w", err)
	}

	fmt.Printf("Successfully updated DAP to %s!\n", release.TagName)
	return nil
}

func getExpectedAssetName() string {
	switch runtime.GOOS {
	case "windows":
		if runtime.GOARCH == "amd64" {
			return "dap-windows-amd64.zip"
		}
	case "linux":
		if runtime.GOARCH == "amd64" {
			return "dap-linux-amd64.tar.gz"
		}
	case "darwin":
		if runtime.GOARCH == "amd64" {
			return "dap-darwin-amd64.tar.gz"
		} else if runtime.GOARCH == "arm64" {
			return "dap-darwin-arm64.tar.gz"
		}
	}
	return ""
}

func downloadFile(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed with status %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func extractFromZip(data []byte, targetFileName string) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, err
	}

	for _, file := range reader.File {
		if filepath.Base(file.Name) == targetFileName {
			return extractZipFile(file)
		}
	}

	return nil, fmt.Errorf("file %q not found in zip archive", targetFileName)
}

func extractZipFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func extractFromTarGz(data []byte, targetFileName string) ([]byte, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if filepath.Base(header.Name) == targetFileName && (header.Typeflag == tar.TypeReg || header.Typeflag == tar.TypeRegA) {
			return io.ReadAll(tarReader)
		}
	}

	return nil, fmt.Errorf("file %q not found in tar.gz archive", targetFileName)
}

func replaceExecutable(targetPath string, newBytes []byte) error {
	dir := filepath.Dir(targetPath)

	// Write new binary to a temp file in the same directory first
	tempFile, err := os.CreateTemp(dir, "dap_update_*")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()

	if _, err := tempFile.Write(newBytes); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return err
	}
	if err := tempFile.Chmod(0755); err != nil {
		tempFile.Close()
		os.Remove(tempPath)
		return err
	}
	if err := tempFile.Close(); err != nil {
		os.Remove(tempPath)
		return err
	}

	if runtime.GOOS == "windows" {
		// On Windows, rename existing running exe to .old before moving new one into place
		oldPath := targetPath + fmt.Sprintf(".old_%d", time.Now().Unix())
		if err := os.Rename(targetPath, oldPath); err != nil {
			os.Remove(tempPath)
			return fmt.Errorf("unable to backup existing binary: %w", err)
		}

		if err := os.Rename(tempPath, targetPath); err != nil {
			_ = os.Rename(oldPath, targetPath)
			return fmt.Errorf("unable to place new binary: %w", err)
		}

		_ = os.Remove(oldPath)
		return nil
	}

	if err := os.Rename(tempPath, targetPath); err != nil {
		os.Remove(tempPath)
		return err
	}

	return os.Chmod(targetPath, 0755)
}
