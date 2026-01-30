package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type SubtitleLine struct {
	Index      int
	StartTime  string
	EndTime    string
	Text       string
	Style      string
	OriginalID int
	// ASS-specific fields
	Layer    int
	MarginL  int
	MarginR  int
	MarginV  int
	Effect   string
	RawEvent string // Store original event line for reconstruction
}

type SubtitleFile struct {
	Format    string
	Header    string
	Lines     []SubtitleLine
	LineCount int
	// ASS-specific: store Events section header
	EventsHeader string
}

func ParseFile(path string) (*SubtitleFile, error) {
	if strings.HasSuffix(strings.ToLower(path), ".srt") {
		return parseSRT(path)
	}
	return parseASS(path)
}

// parseASS parses Advanced SubStation Alpha (.ass) subtitle files
func parseASS(path string) (*SubtitleFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open ASS file: %w", err)
	}
	defer file.Close()

	sf := &SubtitleFile{Format: "ass"}
	scanner := bufio.NewScanner(file)

	var headerBuilder strings.Builder
	var inEventsSection bool
	var eventsFormat []string
	lineIndex := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Detect sections
		if strings.HasPrefix(line, "[Events]") {
			inEventsSection = true
			headerBuilder.WriteString(line + "\n")
			continue
		} else if strings.HasPrefix(line, "[") && inEventsSection {
			// New section started, events ended
			inEventsSection = false
		}

		if inEventsSection {
			// Parse Format line
			if strings.HasPrefix(line, "Format:") {
				sf.EventsHeader = line
				headerBuilder.WriteString(line + "\n")
				// Parse format to know column positions
				formatPart := strings.TrimPrefix(line, "Format:")
				parts := strings.Split(formatPart, ",")
				for _, p := range parts {
					eventsFormat = append(eventsFormat, strings.TrimSpace(p))
				}
				continue
			}

			// Parse Dialogue lines
			if strings.HasPrefix(line, "Dialogue:") {
				dialoguePart := strings.TrimPrefix(line, "Dialogue:")

				// ASS format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
				// Text field can contain commas, so we split by comma but limit
				parts := strings.SplitN(dialoguePart, ",", 10)
				if len(parts) >= 10 {
					subLine := SubtitleLine{
						Index:      lineIndex,
						StartTime:  strings.TrimSpace(parts[1]),
						EndTime:    strings.TrimSpace(parts[2]),
						Style:      strings.TrimSpace(parts[3]),
						Text:       strings.TrimSpace(parts[9]), // Text is last field
						RawEvent:   line,
						OriginalID: lineIndex,
					}

					// Parse optional numeric fields
					if layer, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
						subLine.Layer = layer
					}
					if marginL, err := strconv.Atoi(strings.TrimSpace(parts[5])); err == nil {
						subLine.MarginL = marginL
					}
					if marginR, err := strconv.Atoi(strings.TrimSpace(parts[6])); err == nil {
						subLine.MarginR = marginR
					}
					if marginV, err := strconv.Atoi(strings.TrimSpace(parts[7])); err == nil {
						subLine.MarginV = marginV
					}
					subLine.Effect = strings.TrimSpace(parts[8])

					sf.Lines = append(sf.Lines, subLine)
					lineIndex++
				}
			}
			// Skip Comment lines
		} else {
			// Part of header (before [Events] or other sections after)
			headerBuilder.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading ASS file: %w", err)
	}

	sf.Header = headerBuilder.String()
	sf.LineCount = len(sf.Lines)

	return sf, nil
}

// parseSRT parses SubRip (.srt) subtitle files
func parseSRT(path string) (*SubtitleFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open SRT file: %w", err)
	}
	defer file.Close()

	sf := &SubtitleFile{Format: "srt"}
	scanner := bufio.NewScanner(file)

	// SRT format:
	// 1
	// 00:00:01,000 --> 00:00:04,000
	// Text line 1
	// Text line 2
	// (blank line)

	var currentLine SubtitleLine
	var textBuilder strings.Builder
	state := 0 // 0=expecting index, 1=expecting timing, 2=expecting text

	timeRegex := regexp.MustCompile(`(\d{2}:\d{2}:\d{2}[,\.]\d{3})\s*-->\s*(\d{2}:\d{2}:\d{2}[,\.]\d{3})`)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		switch state {
		case 0: // Expecting index
			if line == "" {
				continue
			}
			if idx, err := strconv.Atoi(line); err == nil {
				currentLine = SubtitleLine{Index: idx, OriginalID: idx}
				state = 1
			}

		case 1: // Expecting timing
			if matches := timeRegex.FindStringSubmatch(line); len(matches) >= 3 {
				currentLine.StartTime = matches[1]
				currentLine.EndTime = matches[2]
				textBuilder.Reset()
				state = 2
			}

		case 2: // Expecting text (can be multi-line)
			if line == "" {
				// End of this subtitle block
				currentLine.Text = strings.TrimSpace(textBuilder.String())
				if currentLine.Text != "" {
					sf.Lines = append(sf.Lines, currentLine)
				}
				state = 0
			} else {
				if textBuilder.Len() > 0 {
					textBuilder.WriteString("\n")
				}
				textBuilder.WriteString(line)
			}
		}
	}

	// Handle last subtitle if file doesn't end with blank line
	if state == 2 && textBuilder.Len() > 0 {
		currentLine.Text = strings.TrimSpace(textBuilder.String())
		sf.Lines = append(sf.Lines, currentLine)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading SRT file: %w", err)
	}

	sf.LineCount = len(sf.Lines)
	return sf, nil
}

// RemoveHearingImpairedTags removes common hearing impaired annotations from subtitle text
// Patterns: [...], (...), ♪, and speaker labels like "JOHN:", "Dr. Smith:", "- NARRATOR:"
func RemoveHearingImpairedTags(text string) string {
	// Remove bracketed content [Music], [Laughing], etc.
	text = regexp.MustCompile(`\[.*?\]`).ReplaceAllString(text, "")
	// Remove parenthetical content (sighs), (in Japanese), etc.
	text = regexp.MustCompile(`\(.*?\)`).ReplaceAllString(text, "")
	// Remove music symbols
	text = strings.ReplaceAll(text, "♪", "")
	text = strings.ReplaceAll(text, "♫", "")
	// Remove speaker labels at start of line (handles multi-word, ALL CAPS, titles)
	// Patterns: "JOHN:", "Dr. Smith:", "MR. JONES:", "- NARRATOR:", "CAPTAIN HOOK:"
	text = regexp.MustCompile(`(?m)^-?\s*[A-Z][A-Za-z.\s]*:\s*`).ReplaceAllString(text, "")
	// Remove speaker labels that are ALL CAPS anywhere
	text = regexp.MustCompile(`[A-Z]{2,}[A-Z\s]*:\s*`).ReplaceAllString(text, "")
	// Clean up multiple spaces
	text = regexp.MustCompile(`\s{2,}`).ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

func BatchLines(lines []SubtitleLine, size int) [][]SubtitleLine {
	batches := [][]SubtitleLine{}
	for i := 0; i < len(lines); i += size {
		end := i + size
		if end > len(lines) {
			end = len(lines)
		}
		batches = append(batches, lines[i:end])
	}
	return batches
}

// ReassembleASS reconstructs an ASS file from header and translated lines
func ReassembleASS(header string, lines []SubtitleLine) string {
	var sb strings.Builder
	sb.WriteString(header)

	// Write dialogue lines
	for _, line := range lines {
		// Reconstruct Dialogue line
		// Format: Dialogue: Layer,Start,End,Style,Name,MarginL,MarginR,MarginV,Effect,Text
		sb.WriteString(fmt.Sprintf("Dialogue: %d,%s,%s,%s,,%04d,%04d,%04d,%s,%s\n",
			line.Layer,
			line.StartTime,
			line.EndTime,
			line.Style,
			line.MarginL,
			line.MarginR,
			line.MarginV,
			line.Effect,
			line.Text,
		))
	}

	return sb.String()
}

func ReassembleSRT(lines []SubtitleLine) string {
	var sb strings.Builder
	for i, line := range lines {
		sb.WriteString(fmt.Sprintf("%d\n", i+1))
		sb.WriteString(fmt.Sprintf("%s --> %s\n", line.StartTime, line.EndTime))
		sb.WriteString(line.Text)
		sb.WriteString("\n\n")
	}
	return sb.String()
}
