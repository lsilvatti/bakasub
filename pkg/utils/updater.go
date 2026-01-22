package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	PublishedAt string `json:"published_at"`
	HTMLURL     string `json:"html_url"`
}

// MsgUpdateAvailable is sent when a new version is detected
type MsgUpdateAvailable struct {
	CurrentVersion string
	LatestVersion  string
	ReleaseURL     string
}

// MsgUpdateCheckFailed is sent when the update check fails
type MsgUpdateCheckFailed struct {
	Err error
}

// CheckForUpdates queries GitHub API for the latest release
func CheckForUpdates(currentVersion string) tea.Cmd {
	return func() tea.Msg {
		// Add a small delay to not block startup
		time.Sleep(500 * time.Millisecond)

		client := &http.Client{
			Timeout: 5 * time.Second,
		}

		apiURL := "https://api.github.com/repos/lsilvatti/bakasub/releases/latest"

		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			return MsgUpdateCheckFailed{Err: err}
		}

		// Set user agent (GitHub requires this)
		req.Header.Set("User-Agent", "BakaSub/"+currentVersion)

		resp, err := client.Do(req)
		if err != nil {
			return MsgUpdateCheckFailed{Err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return MsgUpdateCheckFailed{
				Err: fmt.Errorf("GitHub API returned status %d", resp.StatusCode),
			}
		}

		var release GitHubRelease
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return MsgUpdateCheckFailed{Err: err}
		}

		// Compare versions
		if isNewerVersion(release.TagName, currentVersion) {
			return MsgUpdateAvailable{
				CurrentVersion: currentVersion,
				LatestVersion:  release.TagName,
				ReleaseURL:     release.HTMLURL,
			}
		}

		// No update available
		return nil
	}
}

// isNewerVersion compares two semantic version strings
// Returns true if remote is newer than local
func isNewerVersion(remote, local string) bool {
	// Strip 'v' prefix if present
	remote = strings.TrimPrefix(remote, "v")
	local = strings.TrimPrefix(local, "v")

	remoteParts := strings.Split(remote, ".")
	localParts := strings.Split(local, ".")

	// Pad to same length
	maxLen := len(remoteParts)
	if len(localParts) > maxLen {
		maxLen = len(localParts)
	}

	for len(remoteParts) < maxLen {
		remoteParts = append(remoteParts, "0")
	}
	for len(localParts) < maxLen {
		localParts = append(localParts, "0")
	}

	// Compare each part
	for i := 0; i < maxLen; i++ {
		var remoteNum, localNum int
		fmt.Sscanf(remoteParts[i], "%d", &remoteNum)
		fmt.Sscanf(localParts[i], "%d", &localNum)

		if remoteNum > localNum {
			return true
		} else if remoteNum < localNum {
			return false
		}
	}

	return false
}

// FormatUpdateNotification returns a styled update notification string
func FormatUpdateNotification(msg MsgUpdateAvailable) string {
	return fmt.Sprintf("[!] UPDATE AVAILABLE: %s → %s",
		msg.CurrentVersion,
		msg.LatestVersion)
}
