// Package utils provides utility functions for BakaSub.
package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// OpenURL opens the given URL in the system's default browser.
// It uses OS-specific commands to launch the browser:
//   - Windows: cmd /c start {url}
//   - Linux:   xdg-open {url}
//   - macOS:   open {url}
//
// Returns an error if the browser command fails to execute.
func OpenURL(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		// Windows: Use 'cmd /c start' to open URL
		// The empty string "" is required as the window title parameter
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		// macOS: Use 'open' command
		cmd = exec.Command("open", url)
	case "linux":
		// Linux: Use 'xdg-open' (standard for most desktop environments)
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Start the command without waiting for it to complete
	// (browser should open in background)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}

	return nil
}

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	return os.UserHomeDir()
}
