package main

import (
	"strconv"
	"strings"
	"unicode"
)

func ProcessText(text string) string {
	tokens := Tokenize(text)

	var result []string

	for _, token := range tokens {
		if isCommand(token) {
			applyCommand(&result, token)
		} else {
			result = append(result, token)
		}
	}

	fixArticles(result)

	return BuildText(result)
}

func isCommand(token string) bool {
	return strings.HasPrefix(token, "(") && strings.HasSuffix(token, ")")
}

func applyCommand(tokens *[]string, cmd string) {
	action, count := parseCommand(cmd)
	if action == "" || count <= 0 {
		return
	}

	indices := lastWordIndices(*tokens, count)
	if len(indices) == 0 {
		return
	}

	switch action {
	case "up":
		for _, idx := range indices {
			(*tokens)[idx] = strings.ToUpper((*tokens)[idx])
		}
	case "low":
		for _, idx := range indices {
			(*tokens)[idx] = strings.ToLower((*tokens)[idx])
		}
	case "cap":
		for _, idx := range indices {
			(*tokens)[idx] = capitalize((*tokens)[idx])
		}
	case "hex":
		idx := indices[len(indices)-1]
		if value, err := strconv.ParseInt((*tokens)[idx], 16, 64); err == nil {
			(*tokens)[idx] = strconv.FormatInt(value, 10)
		}
	case "bin":
		idx := indices[len(indices)-1]
		if value, err := strconv.ParseInt((*tokens)[idx], 2, 64); err == nil {
			(*tokens)[idx] = strconv.FormatInt(value, 10)
		}
	}
}

func parseCommand(cmd string) (string, int) {
	cmd = strings.TrimPrefix(cmd, "(")
	cmd = strings.TrimSuffix(cmd, ")")
	cmd = strings.TrimSpace(cmd)

	parts := strings.Split(cmd, ",")
	action := strings.ToLower(strings.TrimSpace(parts[0]))
	count := 1

	if len(parts) == 2 {
		n, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err == nil && n > 0 {
			count = n
		}
	}

	switch action {
	case "up", "low", "cap", "hex", "bin":
		return action, count
	default:
		return "", 0
	}
}

func lastWordIndices(tokens []string, count int) []int {
	var indices []int

	for i := len(tokens) - 1; i >= 0 && len(indices) < count; i-- {
		if isWord(tokens[i]) {
			indices = append([]int{i}, indices...)
		}
	}

	return indices
}

func isWord(token string) bool {
	if token == "'" {
		return false
	}
	if isCommand(token) {
		return false
	}
	for _, r := range token {
		if isPunctuationRune(r) {
			return false
		}
	}
	return token != ""
}

func capitalize(word string) string {
	runes := []rune(strings.ToLower(word))
	if len(runes) == 0 {
		return word
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func fixArticles(tokens []string) {
	for i := 0; i < len(tokens); i++ {
		if strings.EqualFold(tokens[i], "a") {
			next := nextWord(tokens, i+1)
			if next == "" {
				continue
			}

			first := firstLetter(next)
			if first == 0 {
				continue
			}

			if isVowelOrH(first) {
				if tokens[i] == "A" {
					tokens[i] = "An"
				} else {
					tokens[i] = "an"
				}
			}
		}
	}
}

func nextWord(tokens []string, start int) string {
	for i := start; i < len(tokens); i++ {
		if isWord(tokens[i]) {
			return tokens[i]
		}
	}
	return ""
}

func firstLetter(word string) rune {
	for _, r := range word {
		if unicode.IsLetter(r) {
			return unicode.ToLower(r)
		}
	}
	return 0
}

func isVowelOrH(r rune) bool {
	return r == 'a' || r == 'e' || r == 'i' || r == 'o' || r == 'u' || r == 'h'
}

func BuildText(tokens []string) string {
	var b strings.Builder
	quoteOpen := false

	for i, token := range tokens {
		if i == 0 {
			b.WriteString(token)
			if token == "'" {
				quoteOpen = true
			}
			continue
		}

		prev := tokens[i-1]

		switch {
		case token == "'":
			if quoteOpen {
				// закрывающая кавычка — вплотную
				b.WriteString(token)
				quoteOpen = false
			} else {
				// открывающая кавычка
				if shouldAddSpaceBeforeOpeningQuote(prev) {
					b.WriteString(" ")
				}
				b.WriteString(token)
				quoteOpen = true
			}

		case isPunctuationToken(token):
			// пунктуация всегда вплотную к предыдущему слову
			b.WriteString(token)

		default:
			// слово
			if prev == "'" && quoteOpen {
				// после открывающей кавычки без пробела
				b.WriteString(token)
			} else {
				b.WriteString(" ")
				b.WriteString(token)
			}
		}
	}

	return b.String()
}

func isPunctuationToken(token string) bool {
	if token == "" {
		return false
	}
	for _, r := range token {
		if !isPunctuationRune(r) {
			return false
		}
	}
	return true
}

func shouldAddSpaceBeforeOpeningQuote(prev string) bool {
	if prev == "" {
		return false
	}
	if prev == "'" {
		return false
	}
	return true
}
