package utils

import "testing"

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		remote   string
		local    string
		expected bool
	}{
		// Newer versions
		{"v1.1.0", "v1.0.0", true},
		{"v2.0.0", "v1.9.9", true},
		{"v1.0.1", "v1.0.0", true},
		{"1.1.0", "1.0.0", true}, // Without 'v' prefix

		// Older versions
		{"v1.0.0", "v1.1.0", false},
		{"v1.9.9", "v2.0.0", false},

		// Same versions
		{"v1.0.0", "v1.0.0", false},
		{"1.0.0", "1.0.0", false},

		// Edge cases
		{"v1.0", "v1.0.0", false},
		{"v1.0.0", "v1.0", false},
	}

	for _, tt := range tests {
		result := isNewerVersion(tt.remote, tt.local)
		if result != tt.expected {
			t.Errorf("isNewerVersion(%q, %q) = %v; want %v",
				tt.remote, tt.local, result, tt.expected)
		}
	}
}
