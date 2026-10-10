package output

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type SummaryLine struct {
	Symbol string
	Style  lipgloss.Style
	Count  int
	Label  string
}

func Summary(lines ...SummaryLine) string {
	rendered := make([]string, 0, len(lines))

	for _, line := range lines {
		if line.Count == 0 {
			continue
		}

		rendered = append(rendered, "  "+line.Style.Render(fmt.Sprintf("%s %d %s", line.Symbol, line.Count, line.Label)))
	}

	return "Summary:\n" + strings.Join(rendered, "\n")
}
