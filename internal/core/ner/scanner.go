// Package ner provides Named Entity Recognition scanning for subtitle files.
// It extracts capitalized entities (names, places) for the Volatile Glossary.
package ner

import (
	"regexp"
	"sort"
	"strings"
	"unicode"

	"github.com/lsilvatti/bakasub/internal/core/parser"
)

// Entity represents a detected named entity
type Entity struct {
	Text       string
	Type       EntityType
	Confidence float64 // 0.0 - 1.0
	Count      int     // How many times it appears
}

// EntityType categorizes the kind of entity
type EntityType string

const (
	EntityName   EntityType = "Name"   // Character names
	EntityPlace  EntityType = "Place"  // Location names
	EntityAttack EntityType = "Attack" // Attack/technique names (anime)
	EntityTitle  EntityType = "Title"  // Titles (honorifics, job titles)
)

// Scanner performs NER on subtitle lines
type Scanner struct {
	// Common words to exclude (articles, prepositions, etc.)
	stopWords map[string]bool
	// Common Japanese honorifics
	honorifics []string
	// Attack name patterns
	attackPatterns []*regexp.Regexp
	// Title patterns
	titlePatterns []*regexp.Regexp
}

// NewScanner creates a new NER scanner with default configuration
func NewScanner() *Scanner {
	s := &Scanner{
		stopWords: make(map[string]bool),
		honorifics: []string{
			"-san", "-kun", "-chan", "-sama", "-sensei", "-senpai", "-dono",
			"-tan", "-han", "-shi",
		},
		attackPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)((?:[A-Z][a-z]+ )+(?:no|style|technique|attack|punch|kick|slash|strike|wave|fist|blade|claw|beam|cannon|burst|impact|smash|break|crush|pierce|cut|destroy|annihilation))`),
			regexp.MustCompile(`[A-Z][a-z]+ [A-Z][a-z]+ no [A-Z][a-z]+`), // "Gomu Gomu no Pistol"
		},
		titlePatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)(captain|commander|lord|lady|king|queen|prince|princess|chief|master|doctor|professor|general|admiral|marshal)`),
		},
	}

	// Common English stop words
	commonStopWords := []string{
		"the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for",
		"of", "with", "by", "from", "as", "is", "was", "are", "were", "been",
		"be", "have", "has", "had", "do", "does", "did", "will", "would",
		"could", "should", "may", "might", "must", "shall", "can", "need",
		"this", "that", "these", "those", "it", "its", "he", "she", "they",
		"we", "you", "i", "me", "my", "your", "his", "her", "their", "our",
		"what", "which", "who", "whom", "where", "when", "why", "how",
		"all", "each", "every", "both", "few", "more", "most", "other",
		"some", "such", "no", "nor", "not", "only", "own", "same", "so",
		"than", "too", "very", "just", "also", "now", "then", "here", "there",
		"if", "else", "while", "because", "although", "though", "unless",
		"until", "before", "after", "above", "below", "between", "into",
		"through", "during", "over", "under", "again", "further", "once",
		// Common verbs in dialogue
		"said", "says", "say", "tell", "told", "think", "know", "see", "look",
		"go", "going", "went", "gone", "come", "came", "take", "took", "make",
		"made", "get", "got", "give", "gave", "let", "want", "like",
	}

	for _, w := range commonStopWords {
		s.stopWords[strings.ToLower(w)] = true
	}

	return s
}

// ScanLines extracts entities from subtitle lines
func (s *Scanner) ScanLines(lines []parser.SubtitleLine) []Entity {
	entityCounts := make(map[string]*Entity)

	for _, line := range lines {
		s.extractFromText(line.Text, entityCounts)
	}

	// Convert map to slice and sort by count
	entities := make([]Entity, 0, len(entityCounts))
	for _, e := range entityCounts {
		// Only include entities that appear multiple times or have high confidence
		if e.Count >= 2 || e.Confidence >= 0.8 {
			entities = append(entities, *e)
		}
	}

	// Sort by count (most frequent first)
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].Count > entities[j].Count
	})

	return entities
}

// extractFromText processes a single line of text
func (s *Scanner) extractFromText(text string, entities map[string]*Entity) {
	// Clean text of ASS tags
	cleanText := s.removeASSTags(text)

	// Extract capitalized words/phrases
	s.extractCapitalizedEntities(cleanText, entities)

	// Extract attack names
	s.extractAttackNames(cleanText, entities)

	// Extract names with honorifics
	s.extractHonorificNames(cleanText, entities)
}

// removeASSTags strips ASS formatting tags from text
func (s *Scanner) removeASSTags(text string) string {
	re := regexp.MustCompile(`\{[^}]*\}`)
	return re.ReplaceAllString(text, "")
}

// extractCapitalizedEntities finds capitalized words that might be names
func (s *Scanner) extractCapitalizedEntities(text string, entities map[string]*Entity) {
	// Pattern for capitalized words (2+ chars starting with uppercase)
	words := strings.Fields(text)

	for i, word := range words {
		// Clean punctuation from word
		cleanWord := s.cleanPunctuation(word)
		if len(cleanWord) < 2 {
			continue
		}

		// Check if word starts with uppercase
		runes := []rune(cleanWord)
		if !unicode.IsUpper(runes[0]) {
			continue
		}

		// Skip if it's a stop word
		if s.stopWords[strings.ToLower(cleanWord)] {
			continue
		}

		// Skip if it's the first word of a sentence (common false positive)
		if i == 0 {
			// Still include if it looks like a proper noun (rare word pattern)
			if !s.looksLikeProperNoun(cleanWord) {
				continue
			}
		}

		// Calculate confidence based on word characteristics
		confidence := s.calculateConfidence(cleanWord, i, len(words))

		// Normalize the key
		key := strings.ToLower(cleanWord)

		if existing, ok := entities[key]; ok {
			existing.Count++
			// Update confidence if higher
			if confidence > existing.Confidence {
				existing.Confidence = confidence
			}
		} else {
			entities[key] = &Entity{
				Text:       cleanWord,
				Type:       EntityName, // Default to name
				Confidence: confidence,
				Count:      1,
			}
		}
	}
}

// extractAttackNames finds attack/technique name patterns
func (s *Scanner) extractAttackNames(text string, entities map[string]*Entity) {
	for _, pattern := range s.attackPatterns {
		matches := pattern.FindAllString(text, -1)
		for _, match := range matches {
			match = strings.TrimSpace(match)
			if len(match) < 3 {
				continue
			}

			key := strings.ToLower(match)
			if existing, ok := entities[key]; ok {
				existing.Count++
				existing.Type = EntityAttack
				existing.Confidence = 0.9
			} else {
				entities[key] = &Entity{
					Text:       match,
					Type:       EntityAttack,
					Confidence: 0.9,
					Count:      1,
				}
			}
		}
	}
}

// extractHonorificNames finds names followed by Japanese honorifics
func (s *Scanner) extractHonorificNames(text string, entities map[string]*Entity) {
	for _, hon := range s.honorifics {
		pattern := regexp.MustCompile(`([A-Z][a-z]+)` + regexp.QuoteMeta(hon))
		matches := pattern.FindAllStringSubmatch(text, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				name := match[1]
				key := strings.ToLower(name)

				if existing, ok := entities[key]; ok {
					existing.Count++
					existing.Confidence = 0.95 // High confidence for honorific names
				} else {
					entities[key] = &Entity{
						Text:       name,
						Type:       EntityName,
						Confidence: 0.95,
						Count:      1,
					}
				}
			}
		}
	}
}

// cleanPunctuation removes leading/trailing punctuation from a word
func (s *Scanner) cleanPunctuation(word string) string {
	runes := []rune(word)
	start := 0
	end := len(runes)

	// Remove leading punctuation
	for start < end && !unicode.IsLetter(runes[start]) && !unicode.IsNumber(runes[start]) {
		start++
	}

	// Remove trailing punctuation
	for end > start && !unicode.IsLetter(runes[end-1]) && !unicode.IsNumber(runes[end-1]) {
		end--
	}

	if start >= end {
		return ""
	}

	return string(runes[start:end])
}

// looksLikeProperNoun checks if a word has characteristics of a proper noun
func (s *Scanner) looksLikeProperNoun(word string) bool {
	// Japanese-style names often have specific patterns
	// Names with "ou", "uu", etc.
	if strings.Contains(word, "ou") || strings.Contains(word, "uu") || strings.Contains(word, "ii") {
		return true
	}

	// Names ending in common Japanese suffixes
	suffixes := []string{"ro", "ko", "mi", "ki", "shi", "ta", "da", "na", "ru", "ya"}
	lowerWord := strings.ToLower(word)
	for _, suffix := range suffixes {
		if strings.HasSuffix(lowerWord, suffix) {
			return true
		}
	}

	// Multi-syllable capitalized words are likely proper nouns
	vowels := 0
	for _, r := range strings.ToLower(word) {
		if r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u' {
			vowels++
		}
	}
	return vowels >= 2
}

// calculateConfidence determines how likely a word is to be a named entity
func (s *Scanner) calculateConfidence(word string, position, totalWords int) float64 {
	confidence := 0.5 // Base confidence

	// Longer words are more likely to be names
	if len(word) >= 4 {
		confidence += 0.1
	}
	if len(word) >= 6 {
		confidence += 0.1
	}

	// Not first word in sentence
	if position > 0 {
		confidence += 0.2
	}

	// Has characteristics of proper noun
	if s.looksLikeProperNoun(word) {
		confidence += 0.2
	}

	// Cap at 1.0
	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}

// MergeWithProjectGlossary merges detected entities with project glossary
// Project glossary entries take precedence
func MergeWithProjectGlossary(entities []Entity, projectGlossary map[string]string) map[string]string {
	merged := make(map[string]string)

	// Add detected entities (preserve original form)
	for _, e := range entities {
		// Only high-confidence entities
		if e.Confidence >= 0.7 {
			merged[e.Text] = e.Text // Preserve in translation
		}
	}

	// Project glossary overrides
	for orig, trans := range projectGlossary {
		merged[orig] = trans
	}

	return merged
}
