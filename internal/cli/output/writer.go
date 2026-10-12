package output

import (
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
)

// Stdout и Stderr убирают цвета, если вывод идёт не в терминал или задан NO_COLOR.
var (
	Stdout = lipgloss.Writer
	Stderr = colorprofile.NewWriter(os.Stderr, os.Environ())
)
