package utils

import (
	"runtime"
	"testing"
)

// TestOpenURLPlatformDetection tests that OpenURL handles the current platform
func TestOpenURLPlatformDetection(t *testing.T) {
	goos := runtime.GOOS

	switch goos {
	case "windows", "darwin", "linux":
		// These are supported platforms - the function should work
		// We can't actually test browser opening, but we can verify no panic
	default:
		// Test that unsupported platform returns error
		err := OpenURL("https://example.com")
		if err == nil {
			t.Error("expected error for unsupported platform")
		}
	}
}

// TestOpenURLEmptyURL tests OpenURL with empty URL
func TestOpenURLEmptyURL(t *testing.T) {
	// Test with invalid URL - should still not panic
	// Might or might not error depending on OS behavior
	err := OpenURL("")
	_ = err
	// We just verify it doesn't panic
}

// TestGetHomeDir tests GetHomeDir function
func TestGetHomeDir(t *testing.T) {
	homeDir, err := GetHomeDir()
	if err != nil {
		t.Fatalf("GetHomeDir failed: %v", err)
	}

	if homeDir == "" {
		t.Error("homeDir should not be empty")
	}
}

// TestGetHomeDirNotEmpty tests that home dir returns a valid path
func TestGetHomeDirNotEmpty(t *testing.T) {
	homeDir, err := GetHomeDir()
	if err != nil {
		t.Fatalf("GetHomeDir failed: %v", err)
	}

	// Home dir should start with /  on unix or contain : on windows
	if runtime.GOOS == "windows" {
		if len(homeDir) < 3 {
			t.Error("Windows home dir should be at least 3 characters (e.g., C:\\)")
		}
	} else {
		if homeDir[0] != '/' {
			t.Errorf("Unix home dir should start with /, got: %q", homeDir)
		}
	}
}

// TestOpenURLDoesNotPanic ensures OpenURL doesn't panic with various inputs
func TestOpenURLDoesNotPanic(t *testing.T) {
	testCases := []struct {
		name string
		url  string
	}{
		{"empty", ""},
		{"valid http", "http://example.com"},
		{"valid https", "https://example.com"},
		{"spaces", "http://example.com/path with spaces"},
		{"unicode", "http://example.com/путь"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("OpenURL panicked with input %q: %v", tc.url, r)
				}
			}()
			// We don't check the error - we just ensure no panic
			_ = OpenURL(tc.url)
		})
	}
}
