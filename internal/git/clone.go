package git

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

// Clone клонирует url в directory. Недостающие родительские каталоги git
// создаёт сам. Запрос пароля в терминале отключён, чтобы команда не зависала.
func Clone(url, directory string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("git", "clone", "--quiet", url, directory)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	err := cmd.Run()

	return strings.TrimSpace(buf.String()), err
}
