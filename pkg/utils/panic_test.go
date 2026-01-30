package utils

import (
	"testing"
)

func TestVersionConstant(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}

	// Version should start with 'v'
	if len(Version) > 0 && Version[0] != 'v' {
		t.Errorf("Version should start with 'v', got %q", Version)
	}
}

func TestRepoURLConstant(t *testing.T) {
	if RepoURL == "" {
		t.Error("RepoURL should not be empty")
	}

	// Should be a valid GitHub URL
	if len(RepoURL) < 10 {
		t.Errorf("RepoURL seems too short: %q", RepoURL)
	}
}

func TestRecoverPanicDoesNotPanic(t *testing.T) {
	// Test that RecoverPanic can be deferred without panicking
	// We can't easily test the actual recovery behavior
	// as it would call os.Exit(1)

	// Just verify the function exists and can be called
	defer func() {
		if r := recover(); r != nil {
			// RecoverPanic itself should not cause a panic
			t.Errorf("RecoverPanic caused a panic: %v", r)
		}
	}()

	// Calling RecoverPanic directly without a panic should be safe
	// (it will just return because recover() returns nil)
	RecoverPanic()
}

func TestRecoverPanicAsDeferred(t *testing.T) {
	// Test that RecoverPanic works as a deferred function
	executed := false

	func() {
		// Note: We can't actually test the full behavior because
		// RecoverPanic calls os.Exit(1) after rendering the BSOD
		// This test just verifies it doesn't cause issues when there's no panic
		defer RecoverPanic()
		executed = true
	}()

	if !executed {
		t.Error("Function body should have executed")
	}
}

func TestBSODStylesExist(t *testing.T) {
	// Test that the styles are defined and usable
	// These are package-level variables, so we're testing they don't cause errors

	// We can't directly access the unexported variables, but we can
	// verify the package compiles correctly with them
}

func TestVersionFormat(t *testing.T) {
	// Version should be in format vX.Y.Z
	if len(Version) < 5 {
		t.Errorf("Version seems too short: %q", Version)
	}

	// Check format: v followed by numbers and dots
	if Version[0] != 'v' {
		t.Errorf("Version should start with 'v': %q", Version)
	}

	// Should contain at least one dot for major.minor
	hasDot := false
	for _, c := range Version {
		if c == '.' {
			hasDot = true
			break
		}
	}

	if !hasDot {
		t.Errorf("Version should contain '.': %q", Version)
	}
}

// TestCenterText tests the centerText helper function
func TestCenterText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		expected int // Expected length
	}{
		{"short text", "Hi", 10, 6}, // 4 spaces + 2 chars
		{"exact width", "12345", 5, 5},
		{"longer than width", "1234567890", 5, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := centerText(tt.text, tt.width)
			if len(result) < len(tt.text) {
				t.Errorf("centerText result shorter than input: got %d, input %d", len(result), len(tt.text))
			}
		})
	}
}

// TestWrapText tests the wrapText helper function
func TestWrapText(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		width  int
		indent string
	}{
		{"empty", "", 10, "  "},
		{"short", "Hello world", 20, "  "},
		{"long", "This is a much longer text that needs to be wrapped across multiple lines", 20, "  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapText(tt.text, tt.width, tt.indent)
			// Just verify no panic and result starts with indent if not empty
			if tt.text != "" && len(result) > 0 && result[:len(tt.indent)] != tt.indent {
				t.Errorf("wrapText result should start with indent")
			}
		})
	}
}

// TestSafeRun tests SafeRun wrapper
func TestSafeRun(t *testing.T) {
	executed := false

	// Test that SafeRun executes the function
	// Note: We can't test panic recovery because it calls os.Exit(1)
	SafeRun(func() {
		executed = true
	})

	if !executed {
		t.Error("SafeRun should execute the provided function")
	}
}
