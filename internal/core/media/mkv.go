package media

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Track represents a single track in an MKV file
type Track struct {
	ID         int                    `json:"id"`
	Type       string                 `json:"type"`       // video, audio, subtitles
	Codec      string                 `json:"codec"`      // h264, aac, ass, etc.
	Language   string                 `json:"language"`   // ISO 639-2 code (eng, jpn, por, etc.)
	Name       string                 `json:"track_name"` // User-defined track name
	Default    bool                   `json:"default_track"`
	Forced     bool                   `json:"forced_track"`
	Properties map[string]interface{} `json:"properties"` // Additional metadata
}

// Attachment represents a file attached to the MKV (fonts, images, etc.)
type Attachment struct {
	ID          int    `json:"id"`
	FileName    string `json:"file_name"`
	MimeType    string `json:"mime_type"`
	Size        int64  `json:"size"`
	Description string `json:"description"`
}

// FileInfo represents complete metadata of an MKV file
type FileInfo struct {
	FileName    string       `json:"file_name"`
	Tracks      []Track      `json:"tracks"`
	Attachments []Attachment `json:"attachments"`
	Container   struct {
		Type     string `json:"type"`
		Duration int64  `json:"duration"` // In milliseconds
	} `json:"container"`
}

// mkvMergeJSON represents the raw JSON output from mkvmerge -J
type mkvMergeJSON struct {
	Container struct {
		Properties struct {
			Duration int64 `json:"duration"`
		} `json:"properties"`
		Type string `json:"type"`
	} `json:"container"`
	Tracks      []mkvTrack      `json:"tracks"`
	Attachments []mkvAttachment `json:"attachments"`
}

type mkvTrack struct {
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Codec      string `json:"codec"`
	Properties struct {
		Language     string `json:"language"`
		TrackName    string `json:"track_name"`
		DefaultTrack bool   `json:"default_track"`
		ForcedTrack  bool   `json:"forced_track"`
	} `json:"properties"`
}

type mkvAttachment struct {
	FileName    string `json:"file_name"`
	ID          int    `json:"id"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	Description string `json:"description"`
}

var (
	// BinPath is the directory containing mkvmerge and mkvextract binaries
	BinPath = "./bin"
)

// SetBinPath sets the directory where MKVToolNix binaries are located
func SetBinPath(path string) {
	BinPath = path
}

// getBinaryPath returns the full path to a binary
func getBinaryPath(name string) string {
	binPath := filepath.Join(BinPath, name)
	if _, err := os.Stat(binPath); err == nil {
		return binPath
	}

	// Fallback to system PATH
	if path, err := exec.LookPath(name); err == nil {
		return path
	}

	return name // Last resort - let exec.Command handle it
}

// Analyze parses an MKV file and returns its metadata
// It wraps `mkvmerge -J {file}` and parses the JSON output
func Analyze(path string) (*FileInfo, error) {
	// Verify file exists
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	// Get mkvmerge binary path
	mkvmerge := getBinaryPath("mkvmerge")

	// Execute mkvmerge -J
	cmd := exec.Command(mkvmerge, "-J", path)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("mkvmerge failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to execute mkvmerge: %w", err)
	}

	// Parse JSON output
	var rawData mkvMergeJSON
	if err := json.Unmarshal(output, &rawData); err != nil {
		return nil, fmt.Errorf("failed to parse mkvmerge JSON: %w", err)
	}

	// Convert to our FileInfo structure
	fileInfo := &FileInfo{
		FileName:    filepath.Base(path),
		Tracks:      make([]Track, 0, len(rawData.Tracks)),
		Attachments: make([]Attachment, 0, len(rawData.Attachments)),
	}

	fileInfo.Container.Type = rawData.Container.Type
	fileInfo.Container.Duration = rawData.Container.Properties.Duration

	// Convert tracks
	for _, t := range rawData.Tracks {
		track := Track{
			ID:         t.ID,
			Type:       t.Type,
			Codec:      t.Codec,
			Language:   t.Properties.Language,
			Name:       t.Properties.TrackName,
			Default:    t.Properties.DefaultTrack,
			Forced:     t.Properties.ForcedTrack,
			Properties: make(map[string]interface{}),
		}
		fileInfo.Tracks = append(fileInfo.Tracks, track)
	}

	// Convert attachments
	for _, a := range rawData.Attachments {
		attachment := Attachment{
			ID:          a.ID,
			FileName:    a.FileName,
			MimeType:    a.ContentType,
			Size:        a.Size,
			Description: a.Description,
		}
		fileInfo.Attachments = append(fileInfo.Attachments, attachment)
	}

	return fileInfo, nil
}

// ExtractTrack extracts a specific track from an MKV file
// It wraps `mkvextract tracks {file} {trackID}:{output}`
func ExtractTrack(inputPath string, trackID int, outputPath string) error {
	// Verify input file exists
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("input file not found: %w", err)
	}

	// Create output directory if needed
	outputDir := filepath.Dir(outputPath)
	if outputDir != "" && outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Get mkvextract binary path
	mkvextract := getBinaryPath("mkvextract")

	// Execute mkvextract
	trackSpec := fmt.Sprintf("%d:%s", trackID, outputPath)
	cmd := exec.Command(mkvextract, "tracks", inputPath, trackSpec)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mkvextract failed: %s - %w", string(output), err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputPath); err != nil {
		return fmt.Errorf("extraction completed but output file not found: %w", err)
	}

	return nil
}

// MuxSource represents a source file or track for muxing
type MuxSource struct {
	FilePath string // Path to source file
	TrackIDs []int  // Specific tracks to include (empty = all)
	NoVideo  bool   // Exclude video tracks
	NoAudio  bool   // Exclude audio tracks
	NoSubs   bool   // Exclude subtitle tracks
	Language string // Set language for all tracks
	Name     string // Set track name
}

// MuxOptions represents options for the muxing operation
type MuxOptions struct {
	Title      string // Container title
	OutputPath string // Output file path
	Sources    []MuxSource
	Quiet      bool // Suppress mkvmerge output
}

// Mux combines multiple sources into a single MKV file
// It wraps `mkvmerge -o {output} [options] {sources}`
func Mux(opts MuxOptions) error {
	if opts.OutputPath == "" {
		return fmt.Errorf("output path is required")
	}

	if len(opts.Sources) == 0 {
		return fmt.Errorf("at least one source is required")
	}

	// Create output directory if needed
	outputDir := filepath.Dir(opts.OutputPath)
	if outputDir != "" && outputDir != "." {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Build mkvmerge command
	mkvmerge := getBinaryPath("mkvmerge")
	args := []string{"-o", opts.OutputPath}

	// Add global options
	if opts.Title != "" {
		args = append(args, "--title", opts.Title)
	}

	if opts.Quiet {
		args = append(args, "--quiet")
	}

	// Add sources
	for i, source := range opts.Sources {
		// Verify source exists
		if _, err := os.Stat(source.FilePath); err != nil {
			return fmt.Errorf("source file not found: %s - %w", source.FilePath, err)
		}

		// Track selection
		if len(source.TrackIDs) > 0 {
			trackList := make([]string, len(source.TrackIDs))
			for j, id := range source.TrackIDs {
				trackList[j] = fmt.Sprintf("%d", id)
			}
			args = append(args, "--audio-tracks", strings.Join(trackList, ","))
			args = append(args, "--video-tracks", strings.Join(trackList, ","))
			args = append(args, "--subtitle-tracks", strings.Join(trackList, ","))
		}

		// Track type exclusions
		if source.NoVideo {
			args = append(args, "--no-video")
		}
		if source.NoAudio {
			args = append(args, "--no-audio")
		}
		if source.NoSubs {
			args = append(args, "--no-subtitles")
		}

		// Track metadata
		if source.Language != "" {
			args = append(args, "--language", fmt.Sprintf("0:%s", source.Language))
		}
		if source.Name != "" {
			args = append(args, "--track-name", fmt.Sprintf("0:%s", source.Name))
		}

		// Add source file
		args = append(args, source.FilePath)

		// Add separator between sources (except last one)
		if i < len(opts.Sources)-1 {
			args = append(args, "+")
		}
	}

	// Execute mkvmerge
	cmd := exec.Command(mkvmerge, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mkvmerge failed: %s - %w", string(output), err)
	}

	// Verify output file was created
	if _, err := os.Stat(opts.OutputPath); err != nil {
		return fmt.Errorf("muxing completed but output file not found: %w", err)
	}

	return nil
}

// DetectLanguageConflict checks if multiple subtitle tracks exist for the same language
// Returns true if a conflict is detected, along with the conflicting track IDs
func DetectLanguageConflict(fileInfo *FileInfo, targetLang string) (bool, []int) {
	// Normalize target language (handle both 2-char and 3-char codes)
	targetLang = strings.ToLower(targetLang)

	// Map common variations
	langMap := map[string][]string{
		"en": {"en", "eng", "en-us", "en-gb"},
		"pt": {"pt", "por", "pt-br", "pt-pt"},
		"es": {"es", "spa", "es-la", "es-es"},
		"ja": {"ja", "jpn", "ja-jp"},
		"fr": {"fr", "fra", "fr-fr"},
		"de": {"de", "deu", "de-de"},
		"it": {"it", "ita", "it-it"},
		"ru": {"ru", "rus", "ru-ru"},
		"zh": {"zh", "chi", "zh-cn", "zh-tw"},
		"ko": {"ko", "kor", "ko-kr"},
	}

	// Build list of acceptable language codes
	acceptableLangs := []string{targetLang}
	for baseLang, variants := range langMap {
		for _, variant := range variants {
			if variant == targetLang {
				acceptableLangs = langMap[baseLang]
				break
			}
		}
	}

	// Find all subtitle tracks matching target language
	var matchingTracks []int
	for _, track := range fileInfo.Tracks {
		if track.Type != "subtitles" {
			continue
		}

		trackLang := strings.ToLower(track.Language)
		for _, acceptLang := range acceptableLangs {
			if trackLang == acceptLang {
				matchingTracks = append(matchingTracks, track.ID)
				break
			}
		}
	}

	// Conflict exists if more than one track matches
	return len(matchingTracks) > 1, matchingTracks
}

// GetSubtitleTracks returns all subtitle tracks from FileInfo
func GetSubtitleTracks(fileInfo *FileInfo) []Track {
	subtitles := make([]Track, 0)
	for _, track := range fileInfo.Tracks {
		if track.Type == "subtitles" {
			subtitles = append(subtitles, track)
		}
	}
	return subtitles
}

// GetTrackByID finds a track by its ID
func GetTrackByID(fileInfo *FileInfo, trackID int) (*Track, error) {
	for _, track := range fileInfo.Tracks {
		if track.ID == trackID {
			return &track, nil
		}
	}
	return nil, fmt.Errorf("track with ID %d not found", trackID)
}

// GetTracksByType returns all tracks of a specific type
func GetTracksByType(fileInfo *FileInfo, trackType string) []Track {
	tracks := make([]Track, 0)
	for _, track := range fileInfo.Tracks {
		if track.Type == trackType {
			tracks = append(tracks, track)
		}
	}
	return tracks
}

// GetTracksByLanguage returns all tracks matching a language
func GetTracksByLanguage(fileInfo *FileInfo, language string) []Track {
	language = strings.ToLower(language)
	tracks := make([]Track, 0)
	for _, track := range fileInfo.Tracks {
		if strings.ToLower(track.Language) == language {
			tracks = append(tracks, track)
		}
	}
	return tracks
}

// HasAttachments checks if the file has any attachments (fonts, etc.)
func HasAttachments(fileInfo *FileInfo) bool {
	return len(fileInfo.Attachments) > 0
}

// ExtractAttachments extracts all attachments from an MKV file to a directory
func ExtractAttachments(inputPath string, outputDir string) error {
	// Verify input file exists
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("input file not found: %w", err)
	}

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get mkvextract binary path
	mkvextract := getBinaryPath("mkvextract")

	// Execute mkvextract attachments
	cmd := exec.Command(mkvextract, "attachments", inputPath, "--destination", outputDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("mkvextract attachments failed: %s - %w", string(output), err)
	}

	return nil
}

// ExtractSubtitleTrack is a convenience wrapper for extracting subtitle tracks
func ExtractSubtitleTrack(inputPath string, trackID int, outputPath string) error {
	return ExtractTrack(inputPath, trackID, outputPath)
}

// MuxSubtitle is a convenience wrapper for muxing a single subtitle back into video
func MuxSubtitle(inputVideo, subtitlePath, outputPath string) error {
	opts := MuxOptions{
		OutputPath: outputPath,
		Sources: []MuxSource{
			{FilePath: inputVideo},
			{FilePath: subtitlePath, Language: "und", Name: "BakaSub AI Translation"},
		},
		Title: "BakaSub AI Translation",
	}
	return Mux(opts)
}
