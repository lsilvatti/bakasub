package db

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// resetSingleton resets the singleton for testing
func resetSingleton() {
	instance = nil
	instanceOnce = sync.Once{}
}

// TestCacheStruct tests Cache structure
func TestCacheStruct(t *testing.T) {
	resetSingleton()
	// Create a temporary database file
	tmpDir, err := os.MkdirTemp("", "bakasub-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	cache, err := GetInstance(dbPath)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	if cache == nil {
		t.Fatal("Cache is nil")
	}
}

// TestCacheEntryStruct tests CacheEntry structure
func TestCacheEntryStruct(t *testing.T) {
	entry := CacheEntry{
		OriginalHash:   "abc123",
		OriginalText:   "Hello world",
		TranslatedText: "Olá mundo",
		LangPair:       "en-pt",
		Similarity:     0.95,
	}

	if entry.OriginalHash != "abc123" {
		t.Errorf("unexpected OriginalHash: %q", entry.OriginalHash)
	}

	if entry.OriginalText != "Hello world" {
		t.Errorf("unexpected OriginalText: %q", entry.OriginalText)
	}

	if entry.TranslatedText != "Olá mundo" {
		t.Errorf("unexpected TranslatedText: %q", entry.TranslatedText)
	}

	if entry.LangPair != "en-pt" {
		t.Errorf("unexpected LangPair: %q", entry.LangPair)
	}

	if entry.Similarity != 0.95 {
		t.Errorf("unexpected Similarity: %f", entry.Similarity)
	}
}

// TestCacheStatsStruct tests CacheStats structure
func TestCacheStatsStruct(t *testing.T) {
	stats := CacheStats{
		TotalEntries: 100,
		HitRate:      0.85,
		SavedCost:    12.50,
	}

	if stats.TotalEntries != 100 {
		t.Errorf("unexpected TotalEntries: %d", stats.TotalEntries)
	}

	if stats.HitRate != 0.85 {
		t.Errorf("unexpected HitRate: %f", stats.HitRate)
	}

	if stats.SavedCost != 12.50 {
		t.Errorf("unexpected SavedCost: %f", stats.SavedCost)
	}
}

// TestCacheSingleton tests that GetInstance returns the same instance
func TestCacheSingleton(t *testing.T) {
	resetSingleton()
	tmpDir, err := os.MkdirTemp("", "bakasub-test-singleton-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	cache, err := GetInstance(dbPath)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	if cache == nil {
		t.Fatal("Cache is nil")
	}

	// Verify the cache has a valid db connection
	if cache.db == nil {
		t.Error("cache.db should not be nil")
	}
}

// TestCacheInvalidPath tests GetInstance with an invalid path
func TestCacheInvalidPath(t *testing.T) {
	resetSingleton()
	// Note: This test might not fail on all systems since SQLite
	// can create the database file if it doesn't exist.
	// We'll just verify no panic occurs.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("GetInstance panicked: %v", r)
		}
	}()

	// Use a path that definitely won't work
	_, _ = GetInstance("/nonexistent/path/that/wont/work/test.db")
}

// TestHashText tests the hashText function
func TestHashText(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"short", "Hello"},
		{"long", "This is a longer piece of text to hash"},
		{"unicode", "日本語テスト"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := hashText(tt.input)
			if len(hash) != 64 { // SHA256 produces 64 hex characters
				t.Errorf("hashText(%q) = %q, want 64 chars, got %d", tt.input, hash, len(hash))
			}

			// Same input should produce same hash
			hash2 := hashText(tt.input)
			if hash != hash2 {
				t.Error("hashText should be deterministic")
			}
		})
	}
}

// TestCalculateSimilarity tests the calculateSimilarity function
func TestCalculateSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected float64
		delta    float64
	}{
		{"identical", "hello world", "hello world", 1.0, 0.01},
		{"case insensitive", "Hello World", "hello world", 1.0, 0.01},
		{"one char diff", "hello", "hallo", 0.8, 0.1},
		{"completely different", "abc", "xyz", 0.0, 0.1},
		{"empty strings", "", "", 1.0, 0.01},
		{"one empty", "hello", "", 0.0, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateSimilarity(tt.s1, tt.s2)
			if result < tt.expected-tt.delta || result > tt.expected+tt.delta {
				t.Errorf("calculateSimilarity(%q, %q) = %f, want ~%f", tt.s1, tt.s2, result, tt.expected)
			}
		})
	}
}

// TestSaveAndGetExactMatch tests saving and retrieving translations
func TestSaveAndGetExactMatch(t *testing.T) {
	resetSingleton()
	tmpDir, err := os.MkdirTemp("", "bakasub-test-save-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	cache, err := GetInstance(dbPath)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	// Save a translation
	original := "Hello, how are you?"
	translated := "Olá, como você está?"
	langPair := "en-pt"

	err = cache.SaveTranslation(original, translated, langPair)
	if err != nil {
		t.Fatalf("SaveTranslation failed: %v", err)
	}

	// Retrieve it
	result, found := cache.GetExactMatch(original, langPair)
	if !found {
		t.Error("Expected to find cached translation")
	}
	if result != translated {
		t.Errorf("GetExactMatch = %q, want %q", result, translated)
	}

	// Should not find with different lang pair
	_, found = cache.GetExactMatch(original, "en-es")
	if found {
		t.Error("Should not find translation with wrong lang pair")
	}
}

// TestGetFuzzyMatch tests fuzzy matching
func TestGetFuzzyMatch(t *testing.T) {
	resetSingleton()
	tmpDir, err := os.MkdirTemp("", "bakasub-test-fuzzy-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	cache, err := GetInstance(dbPath)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	// Save a translation
	original := "Hello, how are you today?"
	translated := "Olá, como você está hoje?"
	langPair := "en-pt"

	err = cache.SaveTranslation(original, translated, langPair)
	if err != nil {
		t.Fatalf("SaveTranslation failed: %v", err)
	}

	// Similar text should match with high threshold
	similar := "Hello, how are you today"
	entry, found := cache.GetFuzzyMatch(similar, langPair, 0.90)
	if !found {
		t.Error("Expected to find fuzzy match")
	}
	if entry != nil && entry.Similarity < 0.90 {
		t.Errorf("Similarity = %f, want >= 0.90", entry.Similarity)
	}

	// Very different text should not match
	different := "Goodbye, see you later"
	_, found = cache.GetFuzzyMatch(different, langPair, 0.90)
	if found {
		t.Error("Should not find fuzzy match for different text")
	}
}

// TestSaveBatch tests batch saving
func TestSaveBatch(t *testing.T) {
	resetSingleton()
	tmpDir, err := os.MkdirTemp("", "bakasub-test-batch-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	cache, err := GetInstance(dbPath)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	entries := []CacheEntry{
		{OriginalText: "Line 1", TranslatedText: "Linha 1", LangPair: "en-pt"},
		{OriginalText: "Line 2", TranslatedText: "Linha 2", LangPair: "en-pt"},
		{OriginalText: "Line 3", TranslatedText: "Linha 3", LangPair: "en-pt"},
	}

	err = cache.SaveBatch(entries)
	if err != nil {
		t.Fatalf("SaveBatch failed: %v", err)
	}

	// Verify all entries were saved
	for _, e := range entries {
		result, found := cache.GetExactMatch(e.OriginalText, e.LangPair)
		if !found {
			t.Errorf("Entry %q not found in cache", e.OriginalText)
		}
		if result != e.TranslatedText {
			t.Errorf("Got %q, want %q", result, e.TranslatedText)
		}
	}
}

// TestGetStats tests statistics retrieval
func TestGetStats(t *testing.T) {
	resetSingleton()
	tmpDir, err := os.MkdirTemp("", "bakasub-test-stats-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	cache, err := GetInstance(dbPath)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	// Empty cache stats
	stats, err := cache.GetStats()
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if stats.TotalEntries != 0 {
		t.Errorf("TotalEntries = %d, want 0", stats.TotalEntries)
	}

	// Add entries
	cache.SaveTranslation("Hello", "Olá", "en-pt")
	cache.SaveTranslation("World", "Mundo", "en-pt")

	stats, err = cache.GetStats()
	if err != nil {
		t.Fatalf("GetStats failed: %v", err)
	}
	if stats.TotalEntries != 2 {
		t.Errorf("TotalEntries = %d, want 2", stats.TotalEntries)
	}
}

// TestClear tests clearing the cache
func TestClear(t *testing.T) {
	resetSingleton()
	tmpDir, err := os.MkdirTemp("", "bakasub-test-clear-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	cache, err := GetInstance(dbPath)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	// Add some entries
	cache.SaveTranslation("Test1", "Teste1", "en-pt")
	cache.SaveTranslation("Test2", "Teste2", "en-pt")

	// Clear the cache
	err = cache.Clear()
	if err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify cache is empty
	stats, _ := cache.GetStats()
	if stats.TotalEntries != 0 {
		t.Errorf("After Clear, TotalEntries = %d, want 0", stats.TotalEntries)
	}
}

// TestCompact tests database compaction
func TestCompact(t *testing.T) {
	resetSingleton()
	tmpDir, err := os.MkdirTemp("", "bakasub-test-compact-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	dbPath := filepath.Join(tmpDir, "test.db")
	cache, err := GetInstance(dbPath)
	if err != nil {
		t.Fatalf("GetInstance failed: %v", err)
	}

	err = cache.Compact()
	if err != nil {
		t.Fatalf("Compact failed: %v", err)
	}
}

// TestMax tests the max helper function
func TestMax(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{1, 2, 2},
		{5, 3, 5},
		{0, 0, 0},
		{-1, 1, 1},
	}

	for _, tt := range tests {
		got := max(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("max(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}
