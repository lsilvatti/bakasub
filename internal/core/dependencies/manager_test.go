package dependencies

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCheck(t *testing.T) {
	// Backup and restore BinDir
	originalBinDir := BinDir
	defer func() { BinDir = originalBinDir }()

	tmpDir := t.TempDir()
	BinDir = tmpDir

	status, err := Check()
	if err != nil {
		t.Fatalf("Check failed: %v", err)
	}

	if status == nil {
		t.Error("status should not be nil")
	}

	// Should check for all dependencies
	expectedBinaries := []string{"ffmpeg", "ffprobe", "mkvmerge", "mkvextract", "mkvpropedit"}
	for _, bin := range expectedBinaries {
		if _, ok := status[bin]; !ok {
			t.Errorf("status should include %q", bin)
		}
	}
}

func TestCheckBinary(t *testing.T) {
	originalBinDir := BinDir
	defer func() { BinDir = originalBinDir }()

	tmpDir := t.TempDir()
	BinDir = tmpDir

	// Test non-existent binary
	if CheckBinary("nonexistent") {
		t.Error("CheckBinary should return false for non-existent binary")
	}

	// Create a fake binary
	binName := "testbin"
	if runtime.GOOS == "windows" {
		binName += ".exe"
	}
	binPath := filepath.Join(tmpDir, binName)
	if err := os.WriteFile(binPath, []byte("fake binary"), 0755); err != nil {
		t.Fatal(err)
	}

	// Test existing binary
	if !CheckBinary("testbin") {
		t.Error("CheckBinary should return true for existing binary")
	}
}

func TestCheckSystemPath(t *testing.T) {
	// Test with common utilities that should exist on most systems
	// Note: This might fail on minimal systems

	// "ls" on Unix, "cmd" on Windows
	var commonCmd string
	if runtime.GOOS == "windows" {
		commonCmd = "cmd"
	} else {
		commonCmd = "ls"
	}

	if !CheckSystemPath(commonCmd) {
		t.Skipf("common command %q not found in PATH, skipping test", commonCmd)
	}

	// Non-existent command
	if CheckSystemPath("nonexistent_command_xyz_12345") {
		t.Error("CheckSystemPath should return false for non-existent command")
	}
}

func TestDependencyStruct(t *testing.T) {
	dep := Dependency{
		Name:           "TestTool",
		WindowsURL:     "https://example.com/windows",
		LinuxURL:       "https://example.com/linux",
		TargetBinaries: []string{"tool1", "tool2"},
	}

	if dep.Name != "TestTool" {
		t.Errorf("unexpected Name: %q", dep.Name)
	}

	if len(dep.TargetBinaries) != 2 {
		t.Errorf("expected 2 target binaries, got %d", len(dep.TargetBinaries))
	}
}

func TestDependencyGetDownloadURL(t *testing.T) {
	dep := Dependency{
		WindowsURL: "https://example.com/windows",
		LinuxURL:   "https://example.com/linux",
	}

	url := dep.GetDownloadURL()

	if runtime.GOOS == "windows" {
		if url != "https://example.com/windows" {
			t.Errorf("expected Windows URL, got %q", url)
		}
	} else {
		if url != "https://example.com/linux" {
			t.Errorf("expected Linux URL, got %q", url)
		}
	}
}

func TestDependenciesList(t *testing.T) {
	if len(Dependencies) != 2 {
		t.Errorf("expected 2 dependencies, got %d", len(Dependencies))
	}

	// Check FFmpeg
	ffmpegFound := false
	mkvFound := false

	for _, dep := range Dependencies {
		if dep.Name == "FFmpeg" {
			ffmpegFound = true
			if len(dep.TargetBinaries) == 0 {
				t.Error("FFmpeg should have target binaries")
			}
		}
		if dep.Name == "MKVToolNix" {
			mkvFound = true
			if len(dep.TargetBinaries) == 0 {
				t.Error("MKVToolNix should have target binaries")
			}
		}
	}

	if !ffmpegFound {
		t.Error("FFmpeg should be in Dependencies list")
	}

	if !mkvFound {
		t.Error("MKVToolNix should be in Dependencies list")
	}
}

func TestProgressReader(t *testing.T) {
	data := []byte("test data for progress reader")

	progressCalled := false
	var lastRead, lastTotal int64

	pr := &progressReader{
		reader:    &fakeReader{data: data},
		totalSize: int64(len(data)),
		progress: func(read, total int64) {
			progressCalled = true
			lastRead = read
			lastTotal = total
		},
	}

	buf := make([]byte, 10)
	n, err := pr.Read(buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if n != 10 {
		t.Errorf("expected to read 10 bytes, got %d", n)
	}

	if !progressCalled {
		t.Error("progress callback should have been called")
	}

	if lastRead != 10 {
		t.Errorf("expected lastRead 10, got %d", lastRead)
	}

	if lastTotal != int64(len(data)) {
		t.Errorf("expected lastTotal %d, got %d", len(data), lastTotal)
	}
}

func TestProgressReaderNilProgress(t *testing.T) {
	data := []byte("test data")

	pr := &progressReader{
		reader:    &fakeReader{data: data},
		totalSize: int64(len(data)),
		progress:  nil, // No progress callback
	}

	buf := make([]byte, 5)
	_, err := pr.Read(buf)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	// Should not panic with nil progress
}

// fakeReader is a simple io.Reader for testing
type fakeReader struct {
	data   []byte
	offset int
}

func (r *fakeReader) Read(p []byte) (int, error) {
	if r.offset >= len(r.data) {
		return 0, nil
	}

	n := copy(p, r.data[r.offset:])
	r.offset += n
	return n, nil
}

func TestBinDirDefault(t *testing.T) {
	if BinDir != "./bin" {
		t.Errorf("expected default BinDir './bin', got %q", BinDir)
	}
}

func TestExtractArchiveTypes(t *testing.T) {
	// Just test that the function handles different extensions
	// without actually having archive files

	tmpDir := t.TempDir()

	// Test with non-existent archive (should fail gracefully)
	err := Extract(filepath.Join(tmpDir, "nonexistent.zip"), []string{"tool"})
	if err == nil {
		t.Error("Extract should fail for non-existent file")
	}
}
