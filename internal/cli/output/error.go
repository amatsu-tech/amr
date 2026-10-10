package output

import (
	"errors"
	"fmt"
	"os"

	"charm.land/lipgloss/v2"
)

func PrintError(err error) {
	var s *silentError
	if errors.As(err, &s) {
		return
	}

	lipgloss.Fprintln(os.Stderr, ErrorStyle.Render("Error: "+err.Error()))

	var h Hinter
	if errors.As(err, &h) {
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, h.Hint())
	}
}
