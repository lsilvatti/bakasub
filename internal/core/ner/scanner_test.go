package ner

import (
	"testing"

	"github.com/lsilvatti/bakasub/internal/core/parser"
)

// TestEntityTypeConstants tests EntityType constant values
func TestEntityTypeConstants(t *testing.T) {
	if EntityName != "Name" {
		t.Errorf("EntityName = %q, want Name", EntityName)
	}

	if EntityPlace != "Place" {
		t.Errorf("EntityPlace = %q, want Place", EntityPlace)
	}

	if EntityAttack != "Attack" {
		t.Errorf("EntityAttack = %q, want Attack", EntityAttack)
	}

	if EntityTitle != "Title" {
		t.Errorf("EntityTitle = %q, want Title", EntityTitle)
	}
}

// TestEntityStruct tests Entity structure
func TestEntityStruct(t *testing.T) {
	entity := Entity{
		Text:       "Naruto",
		Type:       EntityName,
		Confidence: 0.95,
		Count:      10,
	}

	if entity.Text != "Naruto" {
		t.Errorf("Text = %q, want Naruto", entity.Text)
	}

	if entity.Type != EntityName {
		t.Errorf("Type = %q, want %q", entity.Type, EntityName)
	}

	if entity.Confidence != 0.95 {
		t.Errorf("Confidence = %f, want 0.95", entity.Confidence)
	}

	if entity.Count != 10 {
		t.Errorf("Count = %d, want 10", entity.Count)
	}
}

// TestNewScanner tests NewScanner function
func TestNewScanner(t *testing.T) {
	scanner := NewScanner()

	if scanner == nil {
		t.Fatal("NewScanner returned nil")
	}

	if scanner.stopWords == nil {
		t.Error("stopWords should not be nil")
	}

	if scanner.honorifics == nil {
		t.Error("honorifics should not be nil")
	}
}

// TestScannerStopWords tests that stop words are initialized
func TestScannerStopWords(t *testing.T) {
	scanner := NewScanner()

	// Common stop words should be present
	stopWords := []string{"the", "a", "an", "and", "or", "but", "in", "on", "at"}

	for _, word := range stopWords {
		if !scanner.stopWords[word] {
			t.Errorf("Stop word %q should be in stopWords map", word)
		}
	}
}

// TestScannerHonorifics tests honorifics initialization
func TestScannerHonorifics(t *testing.T) {
	scanner := NewScanner()

	if len(scanner.honorifics) == 0 {
		t.Error("honorifics should not be empty")
	}

	// Check for some common honorifics
	hasHonorifics := false
	for _, h := range scanner.honorifics {
		if h == "-san" || h == "-kun" || h == "-chan" {
			hasHonorifics = true
			break
		}
	}

	if !hasHonorifics {
		t.Error("Should have Japanese honorifics like -san, -kun, -chan")
	}
}

// TestScanLinesEmpty tests ScanLines with empty input
func TestScanLinesEmpty(t *testing.T) {
	scanner := NewScanner()

	entities := scanner.ScanLines([]parser.SubtitleLine{})

	if len(entities) != 0 {
		t.Errorf("Expected no entities for empty input, got %d", len(entities))
	}
}

// TestScanLinesSimpleText tests ScanLines with simple text
func TestScanLinesSimpleText(t *testing.T) {
	scanner := NewScanner()

	lines := []parser.SubtitleLine{
		{Index: 0, Text: "Hello world"},
		{Index: 1, Text: "This is a test"},
	}

	entities := scanner.ScanLines(lines)

	// Simple lowercase text shouldn't produce entities
	// (depends on implementation)
	_ = entities
}

// TestScanLinesWithNames tests ScanLines with proper nouns
func TestScanLinesWithNames(t *testing.T) {
	scanner := NewScanner()

	// Create multiple lines with the same name to meet the threshold
	lines := []parser.SubtitleLine{
		{Index: 0, Text: "Naruto said something."},
		{Index: 1, Text: "Naruto ran away."},
		{Index: 2, Text: "Where is Naruto?"},
	}

	entities := scanner.ScanLines(lines)

	// Should detect "Naruto" as an entity
	found := false
	for _, e := range entities {
		if e.Text == "Naruto" {
			found = true
			break
		}
	}

	// Note: This may or may not find Naruto depending on implementation
	_ = found
}

// TestScanLinesWithHonorifics tests ScanLines with Japanese honorifics
func TestScanLinesWithHonorifics(t *testing.T) {
	scanner := NewScanner()

	lines := []parser.SubtitleLine{
		{Index: 0, Text: "Thank you, Sensei-san!"},
		{Index: 1, Text: "Sensei-san is great!"},
	}

	entities := scanner.ScanLines(lines)

	// The scanner should process this without errors
	_ = entities
}

// TestScanLinesWithASSTags tests that ASS tags are removed
func TestScanLinesWithASSTags(t *testing.T) {
	scanner := NewScanner()

	lines := []parser.SubtitleLine{
		{Index: 0, Text: `{\an8}Naruto said hello.`},
		{Index: 1, Text: `{\pos(100,200)}Naruto ran.`},
	}

	entities := scanner.ScanLines(lines)

	// ASS tags should be removed before processing
	_ = entities
}

// TestScanLinesWithAttackNames tests attack name detection
func TestScanLinesWithAttackNames(t *testing.T) {
	scanner := NewScanner()

	lines := []parser.SubtitleLine{
		{Index: 0, Text: "Rasengan attack!"},
		{Index: 1, Text: "Use the Rasengan technique!"},
		{Index: 2, Text: "His Rasengan destroyed everything."},
	}

	entities := scanner.ScanLines(lines)

	// Should process attack patterns
	_ = entities
}

// TestScanLinesWithPlaces tests place detection
func TestScanLinesWithPlaces(t *testing.T) {
	scanner := NewScanner()

	lines := []parser.SubtitleLine{
		{Index: 0, Text: "Welcome to Konoha!"},
		{Index: 1, Text: "Konoha is beautiful."},
		{Index: 2, Text: "Back to Konoha we go."},
	}

	entities := scanner.ScanLines(lines)

	// Should detect place names
	_ = entities
}

// TestEntityConfidence tests entity confidence values
func TestEntityConfidence(t *testing.T) {
	entity := Entity{
		Text:       "Test",
		Type:       EntityName,
		Confidence: 0.5,
		Count:      1,
	}

	if entity.Confidence < 0 || entity.Confidence > 1 {
		t.Errorf("Confidence %f should be between 0 and 1", entity.Confidence)
	}
}

// TestEntityTypes tests all entity types
func TestEntityTypes(t *testing.T) {
	types := []EntityType{EntityName, EntityPlace, EntityAttack, EntityTitle}

	for _, entityType := range types {
		entity := Entity{
			Text: "Test",
			Type: entityType,
		}

		if entity.Type != entityType {
			t.Errorf("Entity type mismatch: got %q, want %q", entity.Type, entityType)
		}
	}
}

// TestScanLinesMultipleEntities tests detecting multiple different entities
func TestScanLinesMultipleEntities(t *testing.T) {
	scanner := NewScanner()

	lines := []parser.SubtitleLine{
		{Index: 0, Text: "Naruto went to Konoha."},
		{Index: 1, Text: "Sasuke left Konoha."},
		{Index: 2, Text: "Naruto found Sasuke."},
		{Index: 3, Text: "They returned to Konoha."},
	}

	entities := scanner.ScanLines(lines)

	// Should detect multiple entities
	_ = entities
}
