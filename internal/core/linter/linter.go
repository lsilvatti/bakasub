package linter

import (
	"fmt"
	"regexp"
	"strings"
)

type Severity string

const (
	SeverityHigh   Severity = "HIGH"
	SeverityMedium Severity = "MED"
	SeverityLow    Severity = "LOW"
)

type Issue struct {
	LineID      int
	Severity    Severity
	IssueType   string
	Content     string
	Suggestion  string
	AutoFixable bool
}

type Result struct {
	Issues    []Issue
	PassedAll bool
}

// CheckOptions configures the linter behavior
type CheckOptions struct {
	SourceLang string            // Source language ISO code
	TargetLang string            // Target language ISO code
	Glossary   map[string]string // Project glossary for mismatch detection
}

// Check runs all quality checks on translated subtitle lines
func Check(lines []string, opts CheckOptions) Result {
	result := Result{
		Issues:    []Issue{},
		PassedAll: true,
	}

	for i, line := range lines {
		// Check 1: ASS tag validation
		if tagIssue := checkASSTags(i+1, line); tagIssue != nil {
			result.Issues = append(result.Issues, *tagIssue)
			result.PassedAll = false
		}

		// Check 2: Bracket matching
		if bracketIssue := checkBrackets(i+1, line); bracketIssue != nil {
			result.Issues = append(result.Issues, *bracketIssue)
			result.PassedAll = false
		}

		// Check 3: Source language residue (English detection)
		if opts.SourceLang != "" && opts.SourceLang != opts.TargetLang {
			if residueIssue := checkSourceResidue(i+1, line, opts.SourceLang); residueIssue != nil {
				result.Issues = append(result.Issues, *residueIssue)
				result.PassedAll = false
			}
		}

		// Check 4: Excessive punctuation
		if punctIssue := checkPunctuation(i+1, line); punctIssue != nil {
			result.Issues = append(result.Issues, *punctIssue)
			result.PassedAll = false
		}

		// Check 5: Glossary mismatch
		if len(opts.Glossary) > 0 {
			if glossaryIssue := checkGlossaryMismatch(i+1, line, opts.Glossary); glossaryIssue != nil {
				result.Issues = append(result.Issues, *glossaryIssue)
				// Glossary mismatches are warnings, not failures
			}
		}
	}

	return result
}

// CheckSimple runs quality checks with minimal options (backwards compatibility)
func CheckSimple(lines []string, sourceLang string) Result {
	return Check(lines, CheckOptions{SourceLang: sourceLang})
}

// AutoFix attempts to fix all auto-fixable issues
func AutoFix(lines []string, issues []Issue) []string {
	fixed := make([]string, len(lines))
	copy(fixed, lines)

	for _, issue := range issues {
		if !issue.AutoFixable {
			continue
		}

		idx := issue.LineID - 1
		if idx < 0 || idx >= len(fixed) {
			continue
		}

		switch issue.IssueType {
		case "Broken ASS Tags":
			fixed[idx] = fixASSTags(fixed[idx])
		case "Bracket Mismatch":
			fixed[idx] = fixBrackets(fixed[idx])
		case "Excessive Punctuation":
			fixed[idx] = fixPunctuation(fixed[idx])
		}
	}

	return fixed
}

// checkASSTags validates ASS subtitle tags
func checkASSTags(lineID int, text string) *Issue {
	// Pattern: {\tag} or {\tag...}
	tagPattern := regexp.MustCompile(`\{[^}]*`)

	// Find unclosed tags
	if tagPattern.MatchString(text) {
		match := tagPattern.FindString(text)
		if !strings.Contains(match, "}") {
			return &Issue{
				LineID:      lineID,
				Severity:    SeverityHigh,
				IssueType:   "Broken ASS Tags",
				Content:     truncate(text, 50),
				Suggestion:  "Add closing '}' to ASS tags",
				AutoFixable: true,
			}
		}
	}

	return nil
}

// checkBrackets validates bracket matching
func checkBrackets(lineID int, text string) *Issue {
	brackets := map[rune]rune{
		'(': ')',
		'[': ']',
		'{': '}',
	}

	stack := []rune{}
	for _, char := range text {
		if closer, isOpen := brackets[char]; isOpen {
			stack = append(stack, closer)
		} else if char == ')' || char == ']' || char == '}' {
			if len(stack) == 0 || stack[len(stack)-1] != char {
				return &Issue{
					LineID:      lineID,
					Severity:    SeverityMedium,
					IssueType:   "Bracket Mismatch",
					Content:     truncate(text, 50),
					Suggestion:  fmt.Sprintf("Mismatched bracket: %c", char),
					AutoFixable: true,
				}
			}
			stack = stack[:len(stack)-1]
		}
	}

	if len(stack) > 0 {
		return &Issue{
			LineID:      lineID,
			Severity:    SeverityMedium,
			IssueType:   "Bracket Mismatch",
			Content:     truncate(text, 50),
			Suggestion:  "Unclosed brackets detected",
			AutoFixable: true,
		}
	}

	return nil
}

// checkSourceResidue detects English words in non-English translations
func checkSourceResidue(lineID int, text string, targetLang string) *Issue {
	// Simple English word detector (common words)
	englishWords := []string{
		"the", "is", "are", "was", "were", "have", "has", "had",
		"hello", "goodbye", "yes", "no", "what", "where", "when",
	}

	lowerText := strings.ToLower(text)
	for _, word := range englishWords {
		pattern := regexp.MustCompile(`\b` + word + `\b`)
		if pattern.MatchString(lowerText) {
			return &Issue{
				LineID:      lineID,
				Severity:    SeverityMedium,
				IssueType:   "English Residual",
				Content:     truncate(text, 50),
				Suggestion:  fmt.Sprintf("English word detected: '%s'", word),
				AutoFixable: false,
			}
		}
	}

	return nil
}

// checkPunctuation detects excessive punctuation
func checkPunctuation(lineID int, text string) *Issue {
	// Check for 3+ consecutive punctuation marks
	pattern := regexp.MustCompile(`[!?.]{3,}`)
	if pattern.MatchString(text) {
		return &Issue{
			LineID:      lineID,
			Severity:    SeverityLow,
			IssueType:   "Excessive Punctuation",
			Content:     truncate(text, 50),
			Suggestion:  "Reduce repeated punctuation",
			AutoFixable: true,
		}
	}

	return nil
}

// checkGlossaryMismatch checks if glossary terms were translated correctly
func checkGlossaryMismatch(lineID int, translatedText string, glossary map[string]string) *Issue {
	lowerText := strings.ToLower(translatedText)

	for original, expected := range glossary {
		// Check if original term appears but expected translation doesn't
		if strings.Contains(lowerText, strings.ToLower(original)) {
			// Check if expected translation also appears (it should)
			if !strings.Contains(lowerText, strings.ToLower(expected)) {
				return &Issue{
					LineID:    lineID,
					Severity:  SeverityLow,
					IssueType: "Glossary Mismatch",
					Content:   truncate(translatedText, 50),
					Suggestion: fmt.Sprintf("Expected '%s' to be translated as '%s'",
						original, expected),
					AutoFixable: false,
				}
			}
		}
	}

	return nil
}

// fixASSTags attempts to close unclosed ASS tags
func fixASSTags(text string) string {
	// Add closing brace to unclosed tags
	pattern := regexp.MustCompile(`(\{[^}]*)$`)
	return pattern.ReplaceAllString(text, "$1}")
}

// fixBrackets attempts to balance brackets
func fixBrackets(text string) string {
	// Simple strategy: remove unmatched brackets
	brackets := map[rune]rune{'(': ')', '[': ']', '{': '}'}
	stack := []rune{}
	result := []rune{}

	for _, char := range text {
		if closer, isOpen := brackets[char]; isOpen {
			stack = append(stack, closer)
			result = append(result, char)
		} else if char == ')' || char == ']' || char == '}' {
			if len(stack) > 0 && stack[len(stack)-1] == char {
				stack = stack[:len(stack)-1]
				result = append(result, char)
			}
			// Skip unmatched closing brackets
		} else {
			result = append(result, char)
		}
	}

	return string(result)
}

// fixPunctuation reduces excessive punctuation
func fixPunctuation(text string) string {
	// Replace 3+ consecutive punctuation with 2
	pattern := regexp.MustCompile(`([!?.]){3,}`)
	return pattern.ReplaceAllString(text, "$1$1")
}

// truncate limits text length for display
func truncate(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "..."
}
