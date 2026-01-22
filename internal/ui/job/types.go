package job

import (
	"github.com/lsilvatti/bakasub/internal/core/media"
)

// Msg types for job setup flow
type (
	// MsgAnalysisComplete is sent when directory analysis finishes
	MsgAnalysisComplete struct {
		Files   []AnalyzedFile
		Success bool
		Error   error
	}

	// MsgConflictResolved is sent when user selects a track in the resolution modal
	MsgConflictResolved struct {
		FileIndex int
		TrackID   int
	}

	// MsgCostEstimated is sent when cost calculation completes
	MsgCostEstimated struct {
		TotalChars    int
		EstimatedCost float64
		TokenCount    int
	}

	// MsgDryRunComplete is sent when simulation finishes
	MsgDryRunComplete struct {
		CanWrite      bool
		TokenCount    int
		EstimatedCost float64
		Warnings      []string
	}

	// MsgGlossaryLoaded is sent when glossary.json is loaded
	MsgGlossaryLoaded struct {
		Terms map[string]string
		Path  string
	}

	// MsgGlossarySaved is sent when glossary is saved
	MsgGlossarySaved struct {
		Success bool
		Error   error
	}

	// MsgStartJob is sent when user confirms job start
	MsgStartJob struct{}

	// MsgCancelJob is sent when user cancels the setup
	MsgCancelJob struct{}
)

// AnalyzedFile represents a single analyzed MKV file
type AnalyzedFile struct {
	Path            string
	Filename        string
	Tracks          []media.Track
	Attachments     []media.Attachment
	HasConflict     bool
	ConflictTracks  []media.Track // Subtitle tracks matching target language
	SelectedTrackID int           // -1 if not resolved
	SubtitleChars   int           // Character count for cost estimation
}

// JobConfig represents the configuration for a job
type JobConfig struct {
	// Source
	InputPath string
	Files     []AnalyzedFile
	BatchMode bool // true if processing multiple files

	// Extraction
	SourceLang      string
	TargetLang      string
	ExtractFonts    bool
	AudioReference  bool
	AutoDetectTrack bool

	// Translation
	MediaType        string // "anime", "movie", "series", "documentary", "youtube"
	AIModel          string
	Temperature      float64
	GlossaryPath     string
	GlossaryTerms    map[string]string
	RemoveHITags     bool
	ContextualPrompt string

	// Muxing
	MuxMode        string // "replace", "new-file"
	TrackTitle     string
	SetDefault     bool
	SetForced      bool
	BackupOriginal bool

	// Cost
	EstimatedChars  int
	EstimatedTokens int
	EstimatedCost   float64
	ModelPricePerM  float64
}
