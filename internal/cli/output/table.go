package output

import (
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

func NewTable(headers ...string) *table.Table {
	return table.New().
		Headers(headers...).
		BorderTop(false).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		BorderHeader(false).
		BorderColumn(false).
		BorderRow(false).
		StyleFunc(func(row, _ int) lipgloss.Style {
			style := lipgloss.NewStyle().Padding(0, 1)

			if row == table.HeaderRow {
				return style.Bold(true)
			}

			return style
		})
}
