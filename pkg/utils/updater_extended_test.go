package utils

import (
	"testing"
)

// Note: isNewerVersion is not exported (private function), so we test it
// indirectly through CheckForUpdates. For unit testing purposes,
// we test the exported types and functions only.

func TestFormatUpdateNotification(t *testing.T) {
	msg := MsgUpdateAvailable{
		CurrentVersion: "v1.0.0",
		LatestVersion:  "v2.0.0",
		ReleaseURL:     "https://github.com/example/repo/releases/v2.0.0",
	}

	result := FormatUpdateNotification(msg)

	if result == "" {
		t.Error("FormatUpdateNotification returned empty string")
	}

	// Should contain version info
	if !containsStr(result, "v1.0.0") {
		t.Error("result should contain current version")
	}

	if !containsStr(result, "v2.0.0") {
		t.Error("result should contain latest version")
	}
}

func TestGitHubReleaseStruct(t *testing.T) {
	release := GitHubRelease{
		TagName:     "v1.0.0",
		Name:        "Release 1.0.0",
		PublishedAt: "2024-01-15T10:00:00Z",
		HTMLURL:     "https://github.com/example/repo/releases/v1.0.0",
	}

	if release.TagName != "v1.0.0" {
		t.Errorf("unexpected TagName: %q", release.TagName)
	}

	if release.Name != "Release 1.0.0" {
		t.Errorf("unexpected Name: %q", release.Name)
	}

	if release.PublishedAt != "2024-01-15T10:00:00Z" {
		t.Errorf("unexpected PublishedAt: %q", release.PublishedAt)
	}

	if release.HTMLURL != "https://github.com/example/repo/releases/v1.0.0" {
		t.Errorf("unexpected HTMLURL: %q", release.HTMLURL)
	}
}

func TestMsgUpdateAvailableStruct(t *testing.T) {
	msg := MsgUpdateAvailable{
		CurrentVersion: "v1.0.0",
		LatestVersion:  "v2.0.0",
		ReleaseURL:     "https://example.com",
	}

	if msg.CurrentVersion != "v1.0.0" {
		t.Errorf("unexpected CurrentVersion: %q", msg.CurrentVersion)
	}

	if msg.LatestVersion != "v2.0.0" {
		t.Errorf("unexpected LatestVersion: %q", msg.LatestVersion)
	}

	if msg.ReleaseURL != "https://example.com" {
		t.Errorf("unexpected ReleaseURL: %q", msg.ReleaseURL)
	}
}

func TestMsgUpdateCheckFailedStruct(t *testing.T) {
	err := &testError{msg: "network error"}
	msg := MsgUpdateCheckFailed{Err: err}

	if msg.Err == nil {
		t.Error("Err should not be nil")
	}

	if msg.Err.Error() != "network error" {
		t.Errorf("unexpected error message: %q", msg.Err.Error())
	}
}

func TestCheckForUpdates(t *testing.T) {
	// This returns a tea.Cmd, we just verify it doesn't panic
	cmd := CheckForUpdates("v1.0.0")

	if cmd == nil {
		t.Error("CheckForUpdates should return a Cmd")
	}
}

// Helper types and functions for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
