package main

import "unicode"

func Tokenize(text string) []string {
	var tokens []string
	runes := []rune(text)

	for i := 0; i < len(runes); {
		// Пропускаем пробелы
		if unicode.IsSpace(runes[i]) {
			i++
			continue
		}

		// Команды в скобках: (up), (low, 3), (hex) ...
		if runes[i] == '(' {
			start := i
			i++
			for i < len(runes) && runes[i] != ')' {
				i++
			}
			if i < len(runes) {
				i++
			}
			tokens = append(tokens, string(runes[start:i]))
			continue
		}

		// Одинарная кавычка
		if runes[i] == '\'' {
			tokens = append(tokens, "'")
			i++
			continue
		}

		// Пунктуация: .,!?;:
		if isPunctuationRune(runes[i]) {
			start := i
			for i < len(runes) && isPunctuationRune(runes[i]) {
				i++
			}
			tokens = append(tokens, string(runes[start:i]))
			continue
		}

		// Обычное слово
		start := i
		for i < len(runes) &&
			!unicode.IsSpace(runes[i]) &&
			runes[i] != '(' &&
			runes[i] != '\'' &&
			!isPunctuationRune(runes[i]) {
			i++
		}
		tokens = append(tokens, string(runes[start:i]))
	}

	return tokens
}

func isPunctuationRune(r rune) bool {
	return r == '.' || r == ',' || r == '!' || r == '?' || r == ':' || r == ';'
}