package db

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/agnivade/levenshtein"
	_ "modernc.org/sqlite"
)

// Cache represents a thread-safe translation cache backed by SQLite
type Cache struct {
	db   *sql.DB
	mu   sync.RWMutex
	path string
}

// CacheEntry represents a cached translation
type CacheEntry struct {
	OriginalHash   string
	OriginalText   string
	TranslatedText string
	LangPair       string
	Similarity     float64 // Only populated for fuzzy matches
}

// CacheStats represents cache statistics
type CacheStats struct {
	TotalEntries int
	HitRate      float64
	SavedCost    float64 // Estimated USD saved
}

var (
	instance     *Cache
	instanceOnce sync.Once
)

// GetInstance returns the singleton cache instance
func GetInstance(dbPath string) (*Cache, error) {
	var initErr error
	instanceOnce.Do(func() {
		instance, initErr = newCache(dbPath)
	})
	return instance, initErr
}

// Open is a convenience function to open/create the cache database
// If dbPath is empty, uses the default path (bakasub.db in config dir)
func Open(dbPath string) (*Cache, error) {
	if dbPath == "" {
		dbPath = "bakasub.db"
	}
	return newCache(dbPath)
}

// newCache creates a new cache instance
func newCache(dbPath string) (*Cache, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to enable WAL mode: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	cache := &Cache{
		db:   db,
		path: dbPath,
	}

	// Initialize schema
	if err := cache.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return cache, nil
}

// initSchema creates the cache table if it doesn't exist
func (c *Cache) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS cache (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		original_hash TEXT NOT NULL,
		original_text TEXT NOT NULL,
		translated_text TEXT NOT NULL,
		lang_pair TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_used DATETIME DEFAULT CURRENT_TIMESTAMP,
		use_count INTEGER DEFAULT 1,
		UNIQUE(original_hash, lang_pair)
	);

	CREATE INDEX IF NOT EXISTS idx_original_hash ON cache(original_hash);
	CREATE INDEX IF NOT EXISTS idx_lang_pair ON cache(lang_pair);
	CREATE INDEX IF NOT EXISTS idx_original_text ON cache(original_text);
	CREATE INDEX IF NOT EXISTS idx_last_used ON cache(last_used);
	`

	_, err := c.db.Exec(schema)
	return err
}

// hashText generates a SHA256 hash of the text
func hashText(text string) string {
	hash := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// calculateSimilarity returns a similarity score between 0.0 (different) and 1.0 (identical)
func calculateSimilarity(s1, s2 string) float64 {
	// Normalize strings (lowercase, trim)
	s1 = strings.ToLower(strings.TrimSpace(s1))
	s2 = strings.ToLower(strings.TrimSpace(s2))

	if s1 == s2 {
		return 1.0
	}

	// Calculate Levenshtein distance
	distance := levenshtein.ComputeDistance(s1, s2)

	// Convert to similarity percentage
	maxLen := max(len(s1), len(s2))
	if maxLen == 0 {
		return 1.0
	}

	similarity := 1.0 - (float64(distance) / float64(maxLen))
	return similarity
}

// GetExactMatch retrieves an exact match from cache
func (c *Cache) GetExactMatch(text, langPair string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hash := hashText(text)

	var translated string
	err := c.db.QueryRow(`
		SELECT translated_text 
		FROM cache 
		WHERE original_hash = ? AND lang_pair = ?
		LIMIT 1
	`, hash, langPair).Scan(&translated)

	if err == sql.ErrNoRows {
		return "", false
	}
	if err != nil {
		return "", false
	}

	// Update usage stats
	go c.updateUsage(hash, langPair)

	return translated, true
}

// GetFuzzyMatch finds the best fuzzy match above the threshold
func (c *Cache) GetFuzzyMatch(text, langPair string, threshold float64) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// First try exact match
	hash := hashText(text)
	var entry CacheEntry
	err := c.db.QueryRow(`
		SELECT original_hash, original_text, translated_text, lang_pair
		FROM cache 
		WHERE original_hash = ? AND lang_pair = ?
		LIMIT 1
	`, hash, langPair).Scan(&entry.OriginalHash, &entry.OriginalText, &entry.TranslatedText, &entry.LangPair)

	if err == nil {
		entry.Similarity = 1.0
		go c.updateUsage(hash, langPair)
		return &entry, true
	}

	// Fuzzy search: get candidates from same language pair
	// Use text length as a heuristic filter - similar length texts are more likely to match
	textLen := len(text)
	minLen := int(float64(textLen) * threshold) // e.g., 95% of original length
	maxLen := int(float64(textLen) / threshold) // e.g., 105% of original length

	rows, err := c.db.Query(`
		SELECT original_hash, original_text, translated_text, lang_pair
		FROM cache 
		WHERE lang_pair = ? AND LENGTH(original_text) BETWEEN ? AND ?
		ORDER BY last_used DESC
		LIMIT 500
	`, langPair, minLen, maxLen)

	if err != nil {
		return nil, false
	}
	defer rows.Close()

	var bestMatch *CacheEntry
	var bestSimilarity float64

	// Calculate similarity for each candidate
	for rows.Next() {
		var candidate CacheEntry
		if err := rows.Scan(&candidate.OriginalHash, &candidate.OriginalText, &candidate.TranslatedText, &candidate.LangPair); err != nil {
			continue
		}

		similarity := calculateSimilarity(text, candidate.OriginalText)
		if similarity >= threshold && similarity > bestSimilarity {
			bestSimilarity = similarity
			bestMatch = &candidate
			bestMatch.Similarity = similarity
		}
	}

	if bestMatch != nil {
		go c.updateUsage(bestMatch.OriginalHash, langPair)
		return bestMatch, true
	}

	return nil, false
}

// SaveTranslation saves a translation to the cache
func (c *Cache) SaveTranslation(original, translated, langPair string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	hash := hashText(original)

	// Check if entry already exists
	var exists int
	err := c.db.QueryRow(`
		SELECT COUNT(*) FROM cache 
		WHERE original_hash = ? AND lang_pair = ?
	`, hash, langPair).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check existing entry: %w", err)
	}

	if exists > 0 {
		// Update existing entry
		_, err = c.db.Exec(`
			UPDATE cache 
			SET translated_text = ?, last_used = CURRENT_TIMESTAMP, use_count = use_count + 1
			WHERE original_hash = ? AND lang_pair = ?
		`, translated, hash, langPair)
		return err
	}

	// Insert new entry
	_, err = c.db.Exec(`
		INSERT INTO cache (original_hash, original_text, translated_text, lang_pair)
		VALUES (?, ?, ?, ?)
	`, hash, original, translated, langPair)

	if err != nil {
		return fmt.Errorf("failed to insert cache entry: %w", err)
	}

	return nil
}

// SaveBatch saves multiple translations in a single transaction
func (c *Cache) SaveBatch(entries []CacheEntry) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO cache (original_hash, original_text, translated_text, lang_pair)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(original_hash, lang_pair) DO UPDATE SET
			translated_text = excluded.translated_text,
			last_used = CURRENT_TIMESTAMP,
			use_count = cache.use_count + 1
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, entry := range entries {
		hash := hashText(entry.OriginalText)
		_, err := stmt.Exec(hash, entry.OriginalText, entry.TranslatedText, entry.LangPair)
		if err != nil {
			return fmt.Errorf("failed to insert entry: %w", err)
		}
	}

	return tx.Commit()
}

// updateUsage updates the last_used timestamp and use_count
func (c *Cache) updateUsage(hash, langPair string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.db.Exec(`
		UPDATE cache 
		SET last_used = CURRENT_TIMESTAMP, use_count = use_count + 1
		WHERE original_hash = ? AND lang_pair = ?
	`, hash, langPair)
}

// GetStats returns cache statistics
func (c *Cache) GetStats() (*CacheStats, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var stats CacheStats

	// Get total entries
	err := c.db.QueryRow("SELECT COUNT(*) FROM cache").Scan(&stats.TotalEntries)
	if err != nil {
		return nil, fmt.Errorf("failed to get total entries: %w", err)
	}

	// Calculate average reuse (hit rate approximation)
	var avgUseCount sql.NullFloat64
	err = c.db.QueryRow("SELECT AVG(use_count) FROM cache").Scan(&avgUseCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get average use count: %w", err)
	}

	if avgUseCount.Valid && avgUseCount.Float64 > 0 {
		stats.HitRate = (avgUseCount.Float64 - 1) / avgUseCount.Float64 * 100
	}

	// Estimate cost savings (assuming $0.15 per 1M tokens, ~300 tokens per entry)
	tokensPerEntry := 300
	costPerMillion := 0.15
	totalTokens := stats.TotalEntries * tokensPerEntry
	stats.SavedCost = (float64(totalTokens) / 1000000) * costPerMillion * (stats.HitRate / 100)

	return &stats, nil
}

// Clear removes all entries from the cache
func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.db.Exec("DELETE FROM cache")
	if err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	// Vacuum to reclaim space
	_, err = c.db.Exec("VACUUM")
	return err
}

// ClearOld removes entries older than the specified days
func (c *Cache) ClearOld(days int) (int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	result, err := c.db.Exec(`
		DELETE FROM cache 
		WHERE last_used < datetime('now', '-' || ? || ' days')
	`, days)

	if err != nil {
		return 0, fmt.Errorf("failed to clear old entries: %w", err)
	}

	return result.RowsAffected()
}

// Close closes the database connection
func (c *Cache) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Compact optimizes the database
func (c *Cache) Compact() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.db.Exec("VACUUM")
	if err != nil {
		return fmt.Errorf("failed to compact database: %w", err)
	}

	return nil
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
