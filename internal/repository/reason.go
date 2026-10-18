package repository

import (
	"fmt"
	"strings"
)

// Длиннее этого причина обрезается, чтобы таблица не разъезжалась.
const maxReasonLength = 160

// gitFailure превращает вывод упавшей git-команды в одну короткую строку для таблицы.
// Полный вывод git бывает на десятки строк (вплоть до git status), и суть в нём теряется.
func gitFailure(output string, err error) string {
	var lines []string

	for line := range strings.Lines(output) {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return err.Error()
	}

	if existing := countContaining(lines, "already exists, no checkout"); existing > 0 {
		if existing == 1 {
			return "restored partly: 1 new file from the stash already exists in the folder, the stash is kept"
		}

		return fmt.Sprintf("restored partly: %d new files from the stash already exist in the folder, the stash is kept", existing)
	}

	if conflicts := countPrefix(lines, "CONFLICT"); conflicts > 0 {
		if conflicts == 1 {
			return "1 file has merge conflicts, the stash is kept"
		}

		return fmt.Sprintf("%d files have merge conflicts, the stash is kept", conflicts)
	}

	for _, line := range lines {
		lower := strings.ToLower(line)

		for _, prefix := range []string{"error:", "fatal:"} {
			if strings.HasPrefix(lower, prefix) {
				return shorten(strings.TrimSpace(line[len(prefix):]))
			}
		}
	}

	return shorten(lines[0])
}

func countContaining(lines []string, text string) int {
	count := 0

	for _, line := range lines {
		if strings.Contains(line, text) {
			count++
		}
	}

	return count
}

func countPrefix(lines []string, prefix string) int {
	count := 0

	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			count++
		}
	}

	return count
}

func shorten(text string) string {
	if len([]rune(text)) <= maxReasonLength {
		return text
	}

	return string([]rune(text)[:maxReasonLength-1]) + "…"
}
