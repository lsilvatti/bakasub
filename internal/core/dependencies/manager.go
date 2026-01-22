package dependencies

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archiver/v3"
)

// Dependency represents an external tool to download
type Dependency struct {
	Name           string
	WindowsURL     string
	LinuxURL       string
	TargetBinaries []string // Binaries to extract from archive
}

// ProgressFunc is called during download to report progress
type ProgressFunc func(bytesRead, totalBytes int64)

var (
	// BinDir is the local directory where binaries are stored
	BinDir = "./bin"

	// Dependencies list of required tools
	Dependencies = []Dependency{
		{
			Name:           "FFmpeg",
			WindowsURL:     "https://ffmpeg.org/download.html#build-windows",
			LinuxURL:       "https://ffmpeg.org/download.html#build-linux",
			TargetBinaries: []string{"ffmpeg", "ffprobe"},
		},
		{
			Name:           "MKVToolNix",
			WindowsURL:     "https://mkvtoolnix.download/downloads.html#windows",
			LinuxURL:       "https://mkvtoolnix.download/downloads.html#linux",
			TargetBinaries: []string{"mkvmerge", "mkvextract", "mkvpropedit"},
		},
	}
)

// GetDownloadURL returns the appropriate download page URL for the current OS
func (d *Dependency) GetDownloadURL() string {
	if runtime.GOOS == "windows" {
		return d.WindowsURL
	}
	return d.LinuxURL
}

// Check verifies if all required binaries exist in BinDir or system PATH// Check verifies if all required binaries exist in BinDir or system PATH
func Check() (map[string]bool, error) {
	status := make(map[string]bool)

	// Ensure bin directory exists
	if err := os.MkdirAll(BinDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create bin directory: %w", err)
	}

	for _, dep := range Dependencies {
		for _, binary := range dep.TargetBinaries {
			binPath := filepath.Join(BinDir, binary)
			if runtime.GOOS == "windows" {
				binPath += ".exe"
			}

			// Check if file exists in BinDir
			info, err := os.Stat(binPath)
			if err == nil && !info.IsDir() {
				status[binary] = true
				continue
			}

			// Also check system PATH
			if CheckSystemPath(binary) {
				status[binary] = true
			} else {
				status[binary] = false
			}
		}
	}

	return status, nil
}

// CheckBinary checks if a specific binary exists
func CheckBinary(name string) bool {
	binPath := filepath.Join(BinDir, name)
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}

	info, err := os.Stat(binPath)
	return err == nil && !info.IsDir()
}

// Download downloads a file from URL with progress reporting
func Download(url, dest string, progress ProgressFunc) error {
	// Create destination file
	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Get content length
	totalBytes := resp.ContentLength

	// Create progress reader
	reader := &progressReader{
		reader:    resp.Body,
		progress:  progress,
		totalSize: totalBytes,
	}

	// Copy with progress
	_, err = io.Copy(out, reader)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// progressReader wraps io.Reader to report progress
type progressReader struct {
	reader    io.Reader
	progress  ProgressFunc
	totalSize int64
	readBytes int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.readBytes += int64(n)

	if pr.progress != nil {
		pr.progress(pr.readBytes, pr.totalSize)
	}

	return n, err
}

// Extract extracts required binaries from archive to BinDir
func Extract(archivePath string, targetBinaries []string) error {
	// Determine archive type from extension
	ext := strings.ToLower(filepath.Ext(archivePath))

	// Create temporary extraction directory
	tempDir := filepath.Join(os.TempDir(), "bakasub-extract")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract archive
	var err error
	switch ext {
	case ".zip", ".7z":
		// Both .zip and .7z can be handled by Unarchive
		err = archiver.Unarchive(archivePath, tempDir)
	case ".xz", ".tar":
		// Handle .tar.xz
		if strings.HasSuffix(archivePath, ".tar.xz") {
			txz := archiver.NewTarXz()
			err = txz.Unarchive(archivePath, tempDir)
		} else {
			err = archiver.Unarchive(archivePath, tempDir)
		}
	case ".appimage":
		// AppImage is executable itself - just copy it
		return extractAppImage(archivePath, targetBinaries)
	default:
		return fmt.Errorf("unsupported archive type: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	// Find and copy target binaries
	found := make(map[string]bool)
	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Check if this is one of our target binaries
		baseName := strings.TrimSuffix(filepath.Base(path), ".exe")
		for _, target := range targetBinaries {
			if baseName == target {
				// Copy to bin directory
				destPath := filepath.Join(BinDir, filepath.Base(path))
				if err := copyFile(path, destPath); err != nil {
					return fmt.Errorf("failed to copy %s: %w", target, err)
				}

				// Make executable on Unix
				if runtime.GOOS != "windows" {
					if err := os.Chmod(destPath, 0755); err != nil {
						return fmt.Errorf("failed to chmod %s: %w", target, err)
					}
				}

				found[target] = true
				break
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// Verify all binaries were found
	for _, target := range targetBinaries {
		if !found[target] {
			return fmt.Errorf("binary not found in archive: %s", target)
		}
	}

	return nil
}

// extractAppImage handles AppImage files (Linux only)
// AppImages can be extracted by running them with --appimage-extract
func extractAppImage(appImagePath string, targetBinaries []string) error {
	// Make the AppImage executable first
	if err := os.Chmod(appImagePath, 0755); err != nil {
		return fmt.Errorf("failed to chmod AppImage: %w", err)
	}

	// Create a temp directory for extraction
	tempDir := filepath.Join(os.TempDir(), "bakasub-appimage-extract")
	os.RemoveAll(tempDir) // Clean any previous extraction
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract the AppImage using --appimage-extract
	cmd := exec.Command(appImagePath, "--appimage-extract")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		// If extraction fails, fall back to creating wrapper scripts
		return extractAppImageFallback(appImagePath, targetBinaries)
	}

	// Find and copy the extracted binaries
	squashfsRoot := filepath.Join(tempDir, "squashfs-root")
	found := make(map[string]bool)

	err := filepath.Walk(squashfsRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if info.IsDir() {
			return nil
		}

		baseName := filepath.Base(path)
		for _, target := range targetBinaries {
			if baseName == target {
				destPath := filepath.Join(BinDir, target)
				if err := copyFile(path, destPath); err != nil {
					return fmt.Errorf("failed to copy %s: %w", target, err)
				}
				if err := os.Chmod(destPath, 0755); err != nil {
					return fmt.Errorf("failed to chmod %s: %w", target, err)
				}
				found[target] = true
				break
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Check all binaries were found
	for _, target := range targetBinaries {
		if !found[target] {
			return fmt.Errorf("binary not found in AppImage: %s", target)
		}
	}

	return nil
}

// extractAppImageFallback creates wrapper scripts that invoke the AppImage
func extractAppImageFallback(appImagePath string, targetBinaries []string) error {
	// Copy the AppImage to bin directory
	appImageDest := filepath.Join(BinDir, "MKVToolNix.AppImage")
	if err := copyFile(appImagePath, appImageDest); err != nil {
		return fmt.Errorf("failed to copy AppImage: %w", err)
	}
	if err := os.Chmod(appImageDest, 0755); err != nil {
		return fmt.Errorf("failed to chmod AppImage: %w", err)
	}

	// Create wrapper scripts for each binary
	for _, target := range targetBinaries {
		wrapperPath := filepath.Join(BinDir, target)
		wrapperContent := fmt.Sprintf("#!/bin/sh\nexec \"%s\" %s \"$@\"\n", appImageDest, target)

		if err := os.WriteFile(wrapperPath, []byte(wrapperContent), 0755); err != nil {
			return fmt.Errorf("failed to create wrapper for %s: %w", target, err)
		}
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// GetBinaryPath returns the full path to a binary in BinDir
func GetBinaryPath(name string) string {
	binPath := filepath.Join(BinDir, name)
	if runtime.GOOS == "windows" {
		binPath += ".exe"
	}
	return binPath
}

// CheckSystemPath checks if a binary is available in system PATH
func CheckSystemPath(name string) bool {
	if runtime.GOOS == "windows" {
		name += ".exe"
	}

	_, err := exec.LookPath(name)
	return err == nil
}

// DownloadAndInstall downloads and installs a dependency
// DEPRECATED: Auto-download removed. Users should install dependencies manually.
func DownloadAndInstall(dep Dependency, progress ProgressFunc) error {
	return fmt.Errorf("%s must be installed manually. Download from: %s", dep.Name, dep.GetDownloadURL())
}

// GetMissingDependencies returns a list of dependencies that need to be installed
func GetMissingDependencies() ([]Dependency, error) {
	status, err := Check()
	if err != nil {
		return nil, err
	}

	missing := make([]Dependency, 0)
	for _, dep := range Dependencies {
		needsInstall := false
		for _, binary := range dep.TargetBinaries {
			if !status[binary] {
				needsInstall = true
				break
			}
		}

		if needsInstall {
			missing = append(missing, dep)
		}
	}

	return missing, nil
}
