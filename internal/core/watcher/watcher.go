package watcher

import (
"context"
"path/filepath"
"strings"
"sync"
"time"

"github.com/fsnotify/fsnotify"
)

// Watcher monitors a directory for new MKV files
type Watcher struct {
watcher    *fsnotify.Watcher
watchPath  string
debounceMap map[string]*time.Timer
mu         sync.Mutex
OnNewFile  func(string) // Callback when new file detected
ctx        context.Context
cancel     context.CancelFunc
}

// New creates a new file watcher
func New(watchPath string) (*Watcher, error) {
fw, err := fsnotify.NewWatcher()
if err != nil {
return nil, err
}

ctx, cancel := context.WithCancel(context.Background())

return &Watcher{
watcher:     fw,
watchPath:   watchPath,
debounceMap: make(map[string]*time.Timer),
ctx:         ctx,
cancel:      cancel,
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
// Log error (could add callback here)
_ = err
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

// Set new debounce timer (wait 2 seconds for file to finish writing)
w.debounceMap[event.Name] = time.AfterFunc(2*time.Second, func() {
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
// Simple check: try to open file exclusively
// In production, might want more robust lock detection
return true // Placeholder
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

// ScanExisting scans for existing files in the directory
func ScanExisting(dir string) ([]string, error) {
matches, err := filepath.Glob(filepath.Join(dir, "*.mkv"))
if err != nil {
return nil, err
}
return matches, nil
}
