package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lsilvatti/bakasub/internal/config"
	"github.com/lsilvatti/bakasub/internal/core/media"
)

// media-demo: Test tool for MKVToolNix wrapper functions
// Usage: go run cmd/media-demo/main.go [path-to-mkv-file]

func main() {
	// Load config to get bin path
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load config, using default bin path\n")
		media.SetBinPath("./bin")
	} else {
		media.SetBinPath(cfg.BinPath)
	}

	// Get input file from args
	if len(os.Args) < 2 {
		fmt.Println("MKVToolNix Wrapper Demo")
		fmt.Println("=======================")
		fmt.Println()
		fmt.Println("Usage: ./bin/media-demo <path-to-mkv-file>")
		fmt.Println()
		fmt.Println("This tool demonstrates the media package functionality:")
		fmt.Println("  1. Analyze - Parse MKV metadata")
		fmt.Println("  2. Detect language conflicts")
		fmt.Println("  3. Extract tracks")
		fmt.Println("  4. Mux operations")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  ./bin/media-demo /path/to/anime.mkv")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Verify file exists
	if _, err := os.Stat(inputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: File not found: %s\n", inputFile)
		os.Exit(1)
	}

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║ MKVToolNix Wrapper Demo                                       ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// 1. Analyze the file
	fmt.Println("📊 Analyzing MKV file...")
	fmt.Printf("   File: %s\n", filepath.Base(inputFile))
	fmt.Println()

	fileInfo, err := media.Analyze(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Analysis failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Analysis complete!")
	fmt.Println()

	// Display container info
	fmt.Println("📦 Container Information:")
	fmt.Printf("   Type: %s\n", fileInfo.Container.Type)
	fmt.Printf("   Duration: %d ms (%.2f minutes)\n",
		fileInfo.Container.Duration,
		float64(fileInfo.Container.Duration)/1000.0/60.0)
	fmt.Println()

	// Display tracks
	fmt.Printf("🎬 Tracks (%d total):\n", len(fileInfo.Tracks))
	fmt.Println("   ┌────┬───────────┬──────┬────────┬─────────────────────────────┐")
	fmt.Println("   │ ID │ Type      │ Lang │ Codec  │ Name/Flags                  │")
	fmt.Println("   ├────┼───────────┼──────┼────────┼─────────────────────────────┤")

	for _, track := range fileInfo.Tracks {
		flags := ""
		if track.Default {
			flags += "[DEFAULT] "
		}
		if track.Forced {
			flags += "[FORCED] "
		}
		if track.Name != "" {
			flags += track.Name
		}

		fmt.Printf("   │ %2d │ %-9s │ %-4s │ %-6s │ %-27s │\n",
			track.ID,
			track.Type,
			track.Language,
			track.Codec,
			truncate(flags, 27),
		)
	}
	fmt.Println("   └────┴───────────┴──────┴────────┴─────────────────────────────┘")
	fmt.Println()

	// Display attachments
	if len(fileInfo.Attachments) > 0 {
		fmt.Printf("📎 Attachments (%d total):\n", len(fileInfo.Attachments))
		fmt.Println("   ┌────┬──────────────────────────────┬──────────────────┬─────────┐")
		fmt.Println("   │ ID │ Filename                     │ MIME Type        │ Size    │")
		fmt.Println("   ├────┼──────────────────────────────┼──────────────────┼─────────┤")

		for _, att := range fileInfo.Attachments {
			fmt.Printf("   │ %2d │ %-28s │ %-16s │ %7s │\n",
				att.ID,
				truncate(att.FileName, 28),
				truncate(att.MimeType, 16),
				formatSize(att.Size),
			)
		}
		fmt.Println("   └────┴──────────────────────────────┴──────────────────┴─────────┘")
		fmt.Println()
	}

	// 2. Check for conflicts
	fmt.Println("🔍 Conflict Detection:")
	testLanguages := []string{"eng", "jpn", "por", "pt-br"}

	for _, lang := range testLanguages {
		hasConflict, trackIDs := media.DetectLanguageConflict(fileInfo, lang)
		if hasConflict {
			fmt.Printf("   ⚠️  CONFLICT detected for '%s': Multiple tracks found (IDs: %v)\n", lang, trackIDs)
		} else if len(trackIDs) > 0 {
			fmt.Printf("   ✅ No conflict for '%s': Single track found (ID: %d)\n", lang, trackIDs[0])
		}
	}
	fmt.Println()

	// 3. Subtitle track analysis
	subtitles := media.GetSubtitleTracks(fileInfo)
	if len(subtitles) > 0 {
		fmt.Printf("💬 Subtitle Tracks (%d found):\n", len(subtitles))
		for _, sub := range subtitles {
			flags := ""
			if sub.Default {
				flags += " [DEFAULT]"
			}
			if sub.Forced {
				flags += " [FORCED]"
			}

			name := sub.Name
			if name == "" {
				name = "(no name)"
			}

			fmt.Printf("   • ID %d: %s (%s) - %s%s\n",
				sub.ID,
				name,
				sub.Language,
				sub.Codec,
				flags,
			)
		}
		fmt.Println()
	}

	// 4. Track extraction demo (optional - only for first subtitle)
	if len(subtitles) > 0 && shouldExtractDemo() {
		fmt.Println("💾 Extraction Demo:")
		firstSub := subtitles[0]

		// Determine extension based on codec
		ext := ".srt"
		if firstSub.Codec == "SubStationAlpha" || firstSub.Codec == "ass" {
			ext = ".ass"
		}

		outputPath := filepath.Join(os.TempDir(), fmt.Sprintf("bakasub_demo_track_%d%s", firstSub.ID, ext))

		fmt.Printf("   Extracting track %d to: %s\n", firstSub.ID, outputPath)

		if err := media.ExtractTrack(inputFile, firstSub.ID, outputPath); err != nil {
			fmt.Printf("   ❌ Extraction failed: %v\n", err)
		} else {
			info, _ := os.Stat(outputPath)
			fmt.Printf("   ✅ Extracted successfully! (%s)\n", formatSize(info.Size()))
			fmt.Printf("   📁 Temp file: %s\n", outputPath)
			fmt.Println("   (This file will be automatically cleaned up)")

			// Clean up
			defer os.Remove(outputPath)
		}
		fmt.Println()
	}

	// 5. Helper functions demo
	fmt.Println("🔧 Helper Functions:")

	videoTracks := media.GetTracksByType(fileInfo, "video")
	audioTracks := media.GetTracksByType(fileInfo, "audio")
	fmt.Printf("   • Video tracks: %d\n", len(videoTracks))
	fmt.Printf("   • Audio tracks: %d\n", len(audioTracks))
	fmt.Printf("   • Subtitle tracks: %d\n", len(subtitles))
	fmt.Printf("   • Has attachments: %v\n", media.HasAttachments(fileInfo))
	fmt.Println()

	fmt.Println("╔═══════════════════════════════════════════════════════════════╗")
	fmt.Println("║ Demo Complete!                                                ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Note: This demo only reads metadata. No files were modified.")
}

// truncate truncates a string to maxLen characters
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// formatSize formats a byte size in human-readable format
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// shouldExtractDemo checks if we should run the extraction demo
// Only extract if user explicitly wants it (to avoid creating temp files unnecessarily)
func shouldExtractDemo() bool {
	// Check for --extract flag
	for _, arg := range os.Args {
		if arg == "--extract" || arg == "-e" {
			return true
		}
	}
	return false
}
