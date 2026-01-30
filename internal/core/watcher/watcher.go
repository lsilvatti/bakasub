package watcher

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// TouchlessConfig configures automatic processing behavior
type TouchlessConfig struct {
	// SubtitleSelection determines how to pick when multiple subs exist
	SubtitleSelection string // "largest", "smallest", "skip"
	// DefaultProfile is the translation profile to use
	DefaultProfile string
	// MuxingStrategy determines output behavior
	MuxingStrategy string // "replace", "new"
	// TargetLang is the default target language
	TargetLang string
}

// Watcher monitors a directory for new MKV files
type Watcher struct {
	watcher       *fsnotify.Watcher
	watchPath     string
	debounceMap   map[string]*time.Timer
	mu            sync.Mutex
	OnNewFile     func(string) // Callback when new file detected
	OnError       func(error)  // Callback for errors
	TouchlessMode bool         // Enable automatic processing
	Touchless     *TouchlessConfig
	ctx           context.Context
	cancel        context.CancelFunc
}

// New creates a new file watcher
func New(watchPath string) (*Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Watcher{
		watcher:       fw,
		watchPath:     watchPath,
		debounceMap:   make(map[string]*time.Timer),
		ctx:           ctx,
		cancel:        cancel,
		TouchlessMode: false,
		Touchless:     &TouchlessConfig{},
	}, nil
}

// Start begins monitoring the directory
func (w *Watcher) Start() error {
	if err := w.watcher.Add(w.watchPath); err != nil {
		return err
	}

	go w.eventLoop()
	return nil
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	w.cancel()
	w.watcher.Close()
}

// eventLoop processes file system events
func (w *Watcher) eventLoop() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			w.handleEvent(event)
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			if w.OnError != nil {
				w.OnError(err)
			}
		}
	}
}

// handleEvent processes a single file event with debouncing
func (w *Watcher) handleEvent(event fsnotify.Event) {
	// Only interested in Create and Write events
	if event.Op&fsnotify.Create != fsnotify.Create &&
		event.Op&fsnotify.Write != fsnotify.Write {
		return
	}

	// Only process .mkv files
	if !strings.HasSuffix(strings.ToLower(event.Name), ".mkv") {
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	// Cancel existing debounce timer
	if timer, exists := w.debounceMap[event.Name]; exists {
		timer.Stop()
	}

	// Set new debounce timer (wait 3 seconds for file to finish writing)
	w.debounceMap[event.Name] = time.AfterFunc(3*time.Second, func() {
		w.processFile(event.Name)
	})
}

// processFile triggers the callback after debounce period
func (w *Watcher) processFile(path string) {
	w.mu.Lock()
	delete(w.debounceMap, path)
	w.mu.Unlock()

	// Check if file is locked (still being written)
	if !w.isFileReady(path) {
		// Retry after another second
		time.AfterFunc(1*time.Second, func() {
			w.processFile(path)
		})
		return
	}

	if w.OnNewFile != nil {
		w.OnNewFile(path)
	}
}

// isFileReady checks if file is no longer locked/being written
func (w *Watcher) isFileReady(path string) bool {
	// Try to open file for exclusive read to verify it's not locked
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	defer file.Close()

	// Get file info
	info, err := file.Stat()
	if err != nil {
		return false
	}

	// Check if file has non-zero size (not empty placeholder)
	if info.Size() == 0 {
		return false
	}

	// Additional check: try to read first byte
	buf := make([]byte, 1)
	_, err = file.Read(buf)
	if err != nil {
		return false
	}

	// Wait a moment and check if size changed (still writing)
	time.Sleep(500 * time.Millisecond)
	info2, err := os.Stat(path)
	if err != nil {
		return false
	}

	// If size changed, file is still being written
	return info.Size() == info2.Size()
}

// WatchDirectory is a convenience function for simple watching
func WatchDirectory(path string, callback func(string)) (*Watcher, error) {
	w, err := New(path)
	if err != nil {
		return nil, err
	}

	w.OnNewFile = callback

	if err := w.Start(); err != nil {
		return nil, err
	}

	return w, nil
}

// WatchDirectoryTouchless starts watcher with touchless configuration
func WatchDirectoryTouchless(path string, config *TouchlessConfig, callback func(string)) (*Watcher, error) {
	w, err := New(path)
	if err != nil {
		return nil, err
	}

	w.TouchlessMode = true
	w.Touchless = config
	w.OnNewFile = callback

	if err := w.Start(); err != nil {
		return nil, err
	}

	return w, nil
}

// ScanExisting scans for existing files in the directory
func ScanExisting(dir string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "*.mkv"))
	if err != nil {
		return nil, err
	}
	return matches, nil
}
