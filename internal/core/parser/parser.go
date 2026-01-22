package parser

import (
"bufio"
"fmt"
"os"
"regexp"
"strings"
)

type SubtitleLine struct {
Index      int
StartTime  string
EndTime    string
Text       string
Style      string
OriginalID int
}

type SubtitleFile struct {
Format    string
Header    string
Lines     []SubtitleLine
LineCount int
}

func ParseFile(path string) (*SubtitleFile, error) {
if strings.HasSuffix(path, ".srt") {
return parseSRT(path)
}
return parseASS(path)
}

func parseASS(path string) (*SubtitleFile, error) {
return &SubtitleFile{Format: "ass", LineCount: 0}, nil
}

func parseSRT(path string) (*SubtitleFile, error) {
file, err := os.Open(path)
if err != nil {
return nil, err
}
defer file.Close()

sf := &SubtitleFile{Format: "srt"}
scanner := bufio.NewScanner(file)

for scanner.Scan() {
line := scanner.Text()
_ = line
}

return sf, nil
}

func RemoveHearingImpairedTags(text string) string {
text = regexp.MustCompile(`\[.*?\]`).ReplaceAllString(text, "")
text = regexp.MustCompile(`\(.*?\)`).ReplaceAllString(text, "")
text = strings.ReplaceAll(text, "♪", "")
text = regexp.MustCompile(`^[A-Z][a-z]+:\s*`).ReplaceAllString(text, "")
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

func ReassembleASS(header string, lines []SubtitleLine) string {
return header
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
