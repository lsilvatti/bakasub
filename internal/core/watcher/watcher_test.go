package watcher

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tmpDir := t.TempDir()

	watcher, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer watcher.Stop()

	if watcher == nil {
		t.Fatal("watcher should not be nil")
	}

	if watcher.watchPath != tmpDir {
		t.Errorf("expected watchPath %q, got %q", tmpDir, watcher.watchPath)
	}

	if watcher.debounceMap == nil {
		t.Error("debounceMap should be initialized")
	}

	if watcher.TouchlessMode {
		t.Error("TouchlessMode should be false by default")
	}
}

func TestWatcherStart(t *testing.T) {
	tmpDir := t.TempDir()

	watcher, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer watcher.Stop()

	err = watcher.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
}

func TestWatcherStop(t *testing.T) {
	tmpDir := t.TempDir()

	watcher, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	watcher.Start()
	watcher.Stop()

	// Should not panic on double stop
	watcher.Stop()
}

func TestWatcherCallback(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping watcher test in short mode")
	}

	tmpDir := t.TempDir()

	watcher, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer watcher.Stop()

	detected := make(chan string, 1)
	watcher.OnNewFile = func(path string) {
		detected <- path
	}

	err = watcher.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Create a .mkv file
	mkvPath := filepath.Join(tmpDir, "test.mkv")
	if err := os.WriteFile(mkvPath, []byte("fake mkv content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Wait for callback (with timeout)
	select {
	case path := <-detected:
		if path != mkvPath {
			t.Errorf("expected path %q, got %q", mkvPath, path)
		}
	case <-time.After(10 * time.Second):
		t.Error("timeout waiting for file detection")
	}
}

func TestWatcherOnlyMKVFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping watcher test in short mode")
	}

	tmpDir := t.TempDir()

	watcher, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer watcher.Stop()

	detected := make(chan string, 1)
	watcher.OnNewFile = func(path string) {
		detected <- path
	}

	err = watcher.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	// Create a non-MKV file
	txtPath := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(txtPath, []byte("text content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Should not detect non-MKV files
	select {
	case path := <-detected:
		t.Errorf("should not detect non-MKV file: %s", path)
	case <-time.After(1 * time.Second):
		// Expected - no detection
	}
}

func TestWatcherErrorCallback(t *testing.T) {
	tmpDir := t.TempDir()

	watcher, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer watcher.Stop()

	errorReceived := make(chan error, 1)
	watcher.OnError = func(err error) {
		errorReceived <- err
	}

	// The error callback is set up but we can't easily trigger an error
	// Just verify it's set
	if watcher.OnError == nil {
		t.Error("OnError should be set")
	}
}

func TestWatchDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	watcher, err := WatchDirectory(tmpDir, func(path string) {
		// Callback for testing
	})

	if err != nil {
		t.Fatalf("WatchDirectory failed: %v", err)
	}
	defer watcher.Stop()

	if watcher.OnNewFile == nil {
		t.Error("OnNewFile callback should be set")
	}
}

func TestWatchDirectoryTouchless(t *testing.T) {
	tmpDir := t.TempDir()

	config := &TouchlessConfig{
		SubtitleSelection: "largest",
		DefaultProfile:    "anime",
		MuxingStrategy:    "replace",
		TargetLang:        "PT-BR",
	}

	watcher, err := WatchDirectoryTouchless(tmpDir, config, func(path string) {})
	if err != nil {
		t.Fatalf("WatchDirectoryTouchless failed: %v", err)
	}
	defer watcher.Stop()

	if !watcher.TouchlessMode {
		t.Error("TouchlessMode should be enabled")
	}

	if watcher.Touchless.SubtitleSelection != "largest" {
		t.Errorf("unexpected SubtitleSelection: %q", watcher.Touchless.SubtitleSelection)
	}

	if watcher.Touchless.DefaultProfile != "anime" {
		t.Errorf("unexpected DefaultProfile: %q", watcher.Touchless.DefaultProfile)
	}

	if watcher.Touchless.TargetLang != "PT-BR" {
		t.Errorf("unexpected TargetLang: %q", watcher.Touchless.TargetLang)
	}
}

func TestScanExisting(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some MKV files
	for i := 0; i < 3; i++ {
		mkvPath := filepath.Join(tmpDir, "video"+string(rune('0'+i))+".mkv")
		if err := os.WriteFile(mkvPath, []byte("fake"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create non-MKV files
	txtPath := filepath.Join(tmpDir, "readme.txt")
	if err := os.WriteFile(txtPath, []byte("text"), 0644); err != nil {
		t.Fatal(err)
	}

	matches, err := ScanExisting(tmpDir)
	if err != nil {
		t.Fatalf("ScanExisting failed: %v", err)
	}

	if len(matches) != 3 {
		t.Errorf("expected 3 MKV files, got %d", len(matches))
	}
}

func TestScanExistingEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	matches, err := ScanExisting(tmpDir)
	if err != nil {
		t.Fatalf("ScanExisting failed: %v", err)
	}

	if len(matches) != 0 {
		t.Errorf("expected 0 files, got %d", len(matches))
	}
}

func TestTouchlessConfigStruct(t *testing.T) {
	config := TouchlessConfig{
		SubtitleSelection: "smallest",
		DefaultProfile:    "movie",
		MuxingStrategy:    "new",
		TargetLang:        "ES",
	}

	if config.SubtitleSelection != "smallest" {
		t.Errorf("unexpected SubtitleSelection: %q", config.SubtitleSelection)
	}

	if config.DefaultProfile != "movie" {
		t.Errorf("unexpected DefaultProfile: %q", config.DefaultProfile)
	}

	if config.MuxingStrategy != "new" {
		t.Errorf("unexpected MuxingStrategy: %q", config.MuxingStrategy)
	}

	if config.TargetLang != "ES" {
		t.Errorf("unexpected TargetLang: %q", config.TargetLang)
	}
}

func TestIsFileReady(t *testing.T) {
	tmpDir := t.TempDir()

	watcher, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer watcher.Stop()

	// Create a complete file
	filePath := filepath.Join(tmpDir, "test.mkv")
	if err := os.WriteFile(filePath, []byte("complete file content"), 0644); err != nil {
		t.Fatal(err)
	}

	ready := watcher.isFileReady(filePath)
	if !ready {
		t.Error("file should be ready")
	}
}

func TestIsFileReadyNonExistent(t *testing.T) {
	tmpDir := t.TempDir()

	watcher, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer watcher.Stop()

	ready := watcher.isFileReady(filepath.Join(tmpDir, "nonexistent.mkv"))
	if ready {
		t.Error("non-existent file should not be ready")
	}
}

func TestIsFileReadyEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()

	watcher, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer watcher.Stop()

	// Create an empty file
	filePath := filepath.Join(tmpDir, "empty.mkv")
	if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	ready := watcher.isFileReady(filePath)
	if ready {
		t.Error("empty file should not be ready")
	}
}
