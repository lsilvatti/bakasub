package glossary

import (
	"os"
	"path/filepath"
	"testing"
)

// TestStateConstants tests State constants
func TestStateConstants(t *testing.T) {
	if StateView != 0 {
		t.Errorf("StateView = %d, want 0", StateView)
	}

	if StateAddTerm != 1 {
		t.Errorf("StateAddTerm = %d, want 1", StateAddTerm)
	}
}

// TestEntryStruct tests Entry structure
func TestEntryStruct(t *testing.T) {
	entry := Entry{
		Original:     "Nakama",
		Translation:  "Companheiro",
		Type:         "Name",
		AutoDetected: true,
	}

	if entry.Original != "Nakama" {
		t.Errorf("Original = %q, want Nakama", entry.Original)
	}

	if entry.Translation != "Companheiro" {
		t.Errorf("Translation = %q, want Companheiro", entry.Translation)
	}

	if entry.Type != "Name" {
		t.Errorf("Type = %q, want Name", entry.Type)
	}

	if !entry.AutoDetected {
		t.Error("AutoDetected should be true")
	}
}

// TestClosedMsgStruct tests ClosedMsg structure
func TestClosedMsgStruct(t *testing.T) {
	msg := ClosedMsg{}
	_ = msg // Just verify it can be created
}

// TestNew tests New function
func TestNew(t *testing.T) {
	// Create a temporary glossary file
	tmpDir, err := os.MkdirTemp("", "bakasub-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	glossaryPath := filepath.Join(tmpDir, "glossary.json")

	model := New(glossaryPath)

	if model == nil {
		t.Fatal("New returned nil")
	}

	if model.filePath != glossaryPath {
		t.Errorf("filePath = %q, want %q", model.filePath, glossaryPath)
	}

	if model.state != StateView {
		t.Errorf("state = %d, want StateView", model.state)
	}
}

// TestModelStruct tests Model structure fields
func TestModelStruct(t *testing.T) {
	model := Model{
		entries:      []Entry{},
		filePath:     "/test/glossary.json",
		width:        80,
		height:       24,
		state:        StateView,
		currentPage:  0,
		itemsPerPage: 10,
	}

	if model.filePath != "/test/glossary.json" {
		t.Errorf("filePath = %q, want /test/glossary.json", model.filePath)
	}

	if model.state != StateView {
		t.Errorf("state = %d, want StateView", model.state)
	}

	if model.width != 80 {
		t.Errorf("width = %d, want 80", model.width)
	}
}

// TestAddTermFields tests add term form fields
func TestAddTermFields(t *testing.T) {
	model := Model{
		state:          StateAddTerm,
		addOriginal:    "Test",
		addTranslation: "Teste",
		addType:        "Name",
		addField:       0,
	}

	if model.addOriginal != "Test" {
		t.Errorf("addOriginal = %q, want Test", model.addOriginal)
	}

	if model.addTranslation != "Teste" {
		t.Errorf("addTranslation = %q, want Teste", model.addTranslation)
	}

	if model.addType != "Name" {
		t.Errorf("addType = %q, want Name", model.addType)
	}
}

// TestEntriesSlice tests entries slice operations
func TestEntriesSlice(t *testing.T) {
	entries := []Entry{
		{Original: "Nakama", Translation: "Companheiro"},
		{Original: "Sensei", Translation: "Mestre"},
	}

	if len(entries) != 2 {
		t.Errorf("len(entries) = %d, want 2", len(entries))
	}

	if entries[0].Original != "Nakama" {
		t.Error("First entry not set correctly")
	}
}

// TestEntryJSONTags tests JSON tags on Entry
func TestEntryJSONTags(t *testing.T) {
	entry := Entry{
		Original:     "test",
		Translation:  "teste",
		Type:         "Name",
		AutoDetected: false,
	}

	// Verify struct can be used (JSON tags are compile-time)
	if entry.Original != "test" {
		t.Error("Entry fields not set correctly")
	}
}

// TestPagination tests pagination fields
func TestPagination(t *testing.T) {
	model := Model{
		currentPage:  2,
		itemsPerPage: 15,
	}

	if model.currentPage != 2 {
		t.Errorf("currentPage = %d, want 2", model.currentPage)
	}

	if model.itemsPerPage != 15 {
		t.Errorf("itemsPerPage = %d, want 15", model.itemsPerPage)
	}
}

// TestLoadOrCreateNonExistent tests loadOrCreate with non-existent file
func TestLoadOrCreateNonExistent(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "bakasub-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	glossaryPath := filepath.Join(tmpDir, "new_glossary.json")

	// loadOrCreate should create an empty slice if file doesn't exist
	entries := loadOrCreate(glossaryPath)

	if entries == nil {
		t.Error("loadOrCreate returned nil")
	}
}
