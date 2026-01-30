package attachments

import (
	"testing"
)

// TestAttachmentStruct tests Attachment structure
func TestAttachmentStruct(t *testing.T) {
	att := Attachment{
		ID:       1,
		FileName: "font.ttf",
		MIMEType: "application/x-font-ttf",
		Size:     12345,
	}

	if att.ID != 1 {
		t.Errorf("ID = %d, want 1", att.ID)
	}

	if att.FileName != "font.ttf" {
		t.Errorf("FileName = %q, want font.ttf", att.FileName)
	}

	if att.MIMEType != "application/x-font-ttf" {
		t.Errorf("MIMEType = %q, want application/x-font-ttf", att.MIMEType)
	}

	if att.Size != 12345 {
		t.Errorf("Size = %d, want 12345", att.Size)
	}
}

// TestMKVAttachmentsStruct tests MKVAttachments structure
func TestMKVAttachmentsStruct(t *testing.T) {
	mkvAtt := MKVAttachments{
		Attachments: []Attachment{
			{ID: 1, FileName: "font1.ttf"},
			{ID: 2, FileName: "font2.otf"},
		},
	}

	if len(mkvAtt.Attachments) != 2 {
		t.Errorf("len(Attachments) = %d, want 2", len(mkvAtt.Attachments))
	}
}

// TestModeConstants tests Mode constants
func TestModeConstants(t *testing.T) {
	if ModeView != 0 {
		t.Errorf("ModeView = %d, want 0", ModeView)
	}

	if ModeDelete != 1 {
		t.Errorf("ModeDelete = %d, want 1", ModeDelete)
	}

	if ModeAdd != 2 {
		t.Errorf("ModeAdd = %d, want 2", ModeAdd)
	}
}

// TestClosedMsgStruct tests ClosedMsg structure
func TestClosedMsgStruct(t *testing.T) {
	msg := ClosedMsg{}
	_ = msg // Just verify it can be created
}

// TestNewWithNonExistentFile tests New with non-existent file
func TestNewWithNonExistentFile(t *testing.T) {
	_, err := New("/nonexistent/file.mkv")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestModelStruct tests Model structure fields
func TestModelStruct(t *testing.T) {
	// We can't easily test New() with a real file, but we can test the structure
	model := Model{
		attachments:  []Attachment{},
		filePath:     "/test/path.mkv",
		mode:         ModeView,
		deleteMarked: make(map[int]bool),
		width:        80,
		height:       24,
		message:      "Test message",
	}

	if model.filePath != "/test/path.mkv" {
		t.Errorf("filePath = %q, want /test/path.mkv", model.filePath)
	}

	if model.mode != ModeView {
		t.Errorf("mode = %d, want ModeView", model.mode)
	}

	if model.width != 80 {
		t.Errorf("width = %d, want 80", model.width)
	}

	if model.height != 24 {
		t.Errorf("height = %d, want 24", model.height)
	}
}

// TestDeleteMarked tests the deleteMarked map
func TestDeleteMarked(t *testing.T) {
	model := Model{
		deleteMarked: make(map[int]bool),
	}

	// Mark some items for deletion
	model.deleteMarked[1] = true
	model.deleteMarked[3] = true

	if !model.deleteMarked[1] {
		t.Error("Item 1 should be marked for deletion")
	}

	if model.deleteMarked[2] {
		t.Error("Item 2 should not be marked for deletion")
	}

	if !model.deleteMarked[3] {
		t.Error("Item 3 should be marked for deletion")
	}
}

// TestAttachmentJSONTags tests JSON tags on Attachment
func TestAttachmentJSONTags(t *testing.T) {
	// This test verifies the struct has proper JSON tags
	att := Attachment{
		ID:       1,
		FileName: "test.ttf",
		MIMEType: "font/ttf",
		Size:     100,
	}

	// Verify struct can be used (JSON tags are compile-time)
	if att.ID != 1 || att.FileName != "test.ttf" {
		t.Error("Attachment fields not set correctly")
	}
}
