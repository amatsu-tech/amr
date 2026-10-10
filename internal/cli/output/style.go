package output

import "charm.land/lipgloss/v2"

var (
	ErrorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	SuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	NoticeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	HeaderStyle  = lipgloss.NewStyle().Bold(true)
)
