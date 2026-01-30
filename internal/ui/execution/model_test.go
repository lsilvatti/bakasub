package execution

import (
	"testing"
	"time"
)

// TestLogLevelConstants tests LogLevel constants
func TestLogLevelConstants(t *testing.T) {
	if LogInfo != 0 {
		t.Errorf("LogInfo = %d, want 0", LogInfo)
	}

	if LogWarn != 1 {
		t.Errorf("LogWarn = %d, want 1", LogWarn)
	}

	if LogError != 2 {
		t.Errorf("LogError = %d, want 2", LogError)
	}

	if LogAI != 3 {
		t.Errorf("LogAI = %d, want 3", LogAI)
	}

	if LogSuccess != 4 {
		t.Errorf("LogSuccess = %d, want 4", LogSuccess)
	}
}

// TestLogEntryStruct tests LogEntry structure
func TestLogEntryStruct(t *testing.T) {
	now := time.Now()
	entry := LogEntry{
		Timestamp: now,
		Level:     LogInfo,
		Message:   "Test message",
	}

	if entry.Timestamp != now {
		t.Error("Timestamp not set correctly")
	}

	if entry.Level != LogInfo {
		t.Errorf("Level = %d, want LogInfo", entry.Level)
	}

	if entry.Message != "Test message" {
		t.Errorf("Message = %q, want Test message", entry.Message)
	}
}

// TestNewLogBuffer tests NewLogBuffer creation
func TestNewLogBuffer(t *testing.T) {
	buffer := NewLogBuffer(100)

	if buffer == nil {
		t.Fatal("NewLogBuffer returned nil")
	}

	if buffer.maxSize != 100 {
		t.Errorf("maxSize = %d, want 100", buffer.maxSize)
	}

	if len(buffer.entries) != 0 {
		t.Errorf("entries should be empty, got %d", len(buffer.entries))
	}
}

// TestLogBufferAddLine tests adding lines to buffer
func TestLogBufferAddLine(t *testing.T) {
	buffer := NewLogBuffer(10)

	buffer.AddLine(LogInfo, "First message")
	buffer.AddLine(LogWarn, "Second message")

	if buffer.Count() != 2 {
		t.Errorf("Count() = %d, want 2", buffer.Count())
	}
}

// TestLogBufferCircular tests circular buffer behavior
func TestLogBufferCircular(t *testing.T) {
	buffer := NewLogBuffer(3)

	buffer.AddLine(LogInfo, "Message 1")
	buffer.AddLine(LogInfo, "Message 2")
	buffer.AddLine(LogInfo, "Message 3")
	buffer.AddLine(LogInfo, "Message 4") // Should remove Message 1

	if buffer.Count() != 3 {
		t.Errorf("Count() = %d, want 3", buffer.Count())
	}

	lines := buffer.GetLines()
	if len(lines) != 3 {
		t.Errorf("GetLines() length = %d, want 3", len(lines))
	}
}

// TestLogBufferGetLines tests GetLines method
func TestLogBufferGetLines(t *testing.T) {
	buffer := NewLogBuffer(10)

	buffer.AddLine(LogInfo, "Test message")
	lines := buffer.GetLines()

	if len(lines) != 1 {
		t.Errorf("GetLines() length = %d, want 1", len(lines))
	}

	if lines[0] == "" {
		t.Error("GetLines() returned empty string")
	}
}

// TestLogBufferGetRawText tests GetRawText method
func TestLogBufferGetRawText(t *testing.T) {
	buffer := NewLogBuffer(10)

	buffer.AddLine(LogInfo, "First message")
	buffer.AddLine(LogWarn, "Second message")

	text := buffer.GetRawText()

	if text == "" {
		t.Error("GetRawText() returned empty string")
	}

	// Should contain newlines
	if len(text) < 10 {
		t.Error("GetRawText() seems too short")
	}
}

// TestLogBufferCount tests Count method
func TestLogBufferCount(t *testing.T) {
	buffer := NewLogBuffer(10)

	if buffer.Count() != 0 {
		t.Errorf("Count() = %d, want 0 for empty buffer", buffer.Count())
	}

	buffer.AddLine(LogInfo, "Test")
	if buffer.Count() != 1 {
		t.Errorf("Count() = %d, want 1", buffer.Count())
	}
}

// TestLogBufferConcurrency tests concurrent access
func TestLogBufferConcurrency(t *testing.T) {
	buffer := NewLogBuffer(100)

	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := 0; i < 50; i++ {
			buffer.AddLine(LogInfo, "Message")
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 50; i++ {
			_ = buffer.GetLines()
			_ = buffer.Count()
		}
		done <- true
	}()

	<-done
	<-done

	// Should not panic - test passes if we get here
}

// TestFormatLogEntry tests log entry formatting
func TestFormatLogEntry(t *testing.T) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     LogInfo,
		Message:   "Test message",
	}

	formatted := FormatLogEntry(entry)

	if formatted == "" {
		t.Error("FormatLogEntry returned empty string")
	}

	// Should contain the message
	if len(formatted) < len("Test message") {
		t.Error("FormatLogEntry seems too short")
	}
}

// TestAllLogLevelsFormat tests formatting all log levels
func TestAllLogLevelsFormat(t *testing.T) {
	levels := []LogLevel{LogInfo, LogWarn, LogError, LogAI, LogSuccess}

	for _, level := range levels {
		entry := LogEntry{
			Timestamp: time.Now(),
			Level:     level,
			Message:   "Test",
		}

		formatted := FormatLogEntry(entry)
		if formatted == "" {
			t.Errorf("FormatLogEntry returned empty for level %d", level)
		}
	}
}

// TestLogBufferClear tests Clear method
func TestLogBufferClear(t *testing.T) {
	buffer := NewLogBuffer(10)

	buffer.AddLine(LogInfo, "Message 1")
	buffer.AddLine(LogInfo, "Message 2")

	if buffer.Count() != 2 {
		t.Errorf("Before Clear: Count() = %d, want 2", buffer.Count())
	}

	buffer.Clear()

	if buffer.Count() != 0 {
		t.Errorf("After Clear: Count() = %d, want 0", buffer.Count())
	}
}

// TestParseLogLine tests ParseLogLine function
func TestParseLogLine(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedLevel LogLevel
	}{
		{"info tag", "[INFO] test message", LogInfo},
		{"info colon", "INFO: test message", LogInfo},
		{"warn tag", "[WARN] test message", LogWarn},
		{"warning tag", "[WARNING] test message", LogWarn},
		{"error tag", "[ERR] test message", LogError},
		{"error colon", "ERROR: test message", LogError},
		{"fail tag", "[FAIL] test message", LogError},
		{"ai tag", "[AI] test message", LogAI},
		{"ai colon", "AI: test message", LogAI},
		{"ok tag", "[OK] test message", LogSuccess},
		{"success tag", "[SUCCESS] test message", LogSuccess},
		{"no tag", "plain message", LogInfo}, // Default to info
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, _ := ParseLogLine(tt.input)
			if level != tt.expectedLevel {
				t.Errorf("ParseLogLine(%q) level = %d, want %d", tt.input, level, tt.expectedLevel)
			}
		})
	}
}

// TestLogBufferMaxSize tests buffer with different max sizes
func TestLogBufferMaxSize(t *testing.T) {
	sizes := []int{1, 5, 100, 1000}

	for _, size := range sizes {
		buffer := NewLogBuffer(size)
		if buffer.maxSize != size {
			t.Errorf("NewLogBuffer(%d) maxSize = %d", size, buffer.maxSize)
		}
	}
}

// TestLogEntryWithEmptyMessage tests log entry with empty message
func TestLogEntryWithEmptyMessage(t *testing.T) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     LogInfo,
		Message:   "",
	}

	formatted := FormatLogEntry(entry)
	if formatted == "" {
		t.Error("FormatLogEntry should return non-empty even with empty message")
	}
}

// TestLogBufferLargeBatch tests adding many entries
func TestLogBufferLargeBatch(t *testing.T) {
	buffer := NewLogBuffer(50)

	for i := 0; i < 100; i++ {
		buffer.AddLine(LogInfo, "Message")
	}

	// Should cap at maxSize
	if buffer.Count() != 50 {
		t.Errorf("Count() = %d, want 50 (maxSize)", buffer.Count())
	}
}
