package lesson

import (
	"strings"
	"unicode"

	"encore.app/learning/content"
)

// EvaluationResult contains the result of evaluating an answer.
type EvaluationResult struct {
	IsCorrect  bool
	Similarity float64
	XPEarned   int
}

// EvaluateAnswer evaluates a user's answer against the correct answer.
func EvaluateAnswer(question *content.Question, userAnswer string, lessonXP int) *EvaluationResult {
	result := &EvaluationResult{}

	switch question.Type {
	case content.QuestionTypeSingleChoice, content.QuestionTypeMultiChoice:
		// Exact match for choice questions
		result.IsCorrect = strings.EqualFold(
			normalizeAnswer(userAnswer),
			normalizeAnswer(question.CorrectAnswer),
		)
		result.Similarity = 1.0
		if !result.IsCorrect {
			result.Similarity = 0.0
		}

	case content.QuestionTypeListenReply:
		// Fuzzy matching for listen & reply questions
		normalizedUser := normalizeAnswer(userAnswer)
		normalizedCorrect := normalizeAnswer(question.CorrectAnswer)

		result.Similarity = calculateSimilarity(normalizedUser, normalizedCorrect)
		result.IsCorrect = result.Similarity >= 0.85 // 85% threshold
	}

	// Calculate XP earned (only for correct answers)
	if result.IsCorrect {
		// Base XP per question (lesson XP divided by expected questions, minimum 1)
		result.XPEarned = max(lessonXP/10, 1)
	}

	return result
}

// normalizeAnswer normalizes an answer for comparison.
func normalizeAnswer(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Remove leading/trailing whitespace
	s = strings.TrimSpace(s)

	// Remove punctuation
	var builder strings.Builder
	for _, r := range s {
		if !unicode.IsPunct(r) || r == '\'' || r == '-' {
			builder.WriteRune(r)
		}
	}
	s = builder.String()

	// Normalize whitespace
	fields := strings.Fields(s)
	s = strings.Join(fields, " ")

	return s
}

// calculateSimilarity calculates the similarity between two strings using Levenshtein distance.
func calculateSimilarity(a, b string) float64 {
	if a == b {
		return 1.0
	}
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}

	distance := levenshteinDistance(a, b)
	maxLen := max(len(a), len(b))

	return 1.0 - float64(distance)/float64(maxLen)
}

// levenshteinDistance calculates the Levenshtein distance between two strings.
func levenshteinDistance(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}

	// Convert to runes for proper Unicode handling
	aRunes := []rune(a)
	bRunes := []rune(b)

	// Create distance matrix
	m := len(aRunes)
	n := len(bRunes)

	// Use two rows instead of full matrix for memory efficiency
	prevRow := make([]int, n+1)
	currRow := make([]int, n+1)

	// Initialize first row
	for j := 0; j <= n; j++ {
		prevRow[j] = j
	}

	// Fill in the rest
	for i := 1; i <= m; i++ {
		currRow[0] = i

		for j := 1; j <= n; j++ {
			cost := 1
			if aRunes[i-1] == bRunes[j-1] {
				cost = 0
			}

			currRow[j] = min(
				prevRow[j]+1,      // deletion
				currRow[j-1]+1,    // insertion
				prevRow[j-1]+cost, // substitution
			)
		}

		// Swap rows
		prevRow, currRow = currRow, prevRow
	}

	return prevRow[n]
}

// GenerateHint generates a hint for an incorrect answer.
func GenerateHint(question *content.Question, userAnswer string) string {
	if question.Hint != "" {
		return question.Hint
	}

	// Generate a basic hint based on the correct answer
	correct := question.CorrectAnswer

	// If the answer is short, hint at the first letter
	if len(correct) <= 10 {
		return "The answer starts with '" + string([]rune(correct)[0]) + "'"
	}

	// For longer answers, hint at the word count
	words := strings.Fields(correct)
	if len(words) > 1 {
		return "The answer has " + string(rune('0'+len(words))) + " words"
	}

	return "Try again!"
}
