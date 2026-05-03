package constant

import (
	"restclient/internal/tui/styles"

	"github.com/charmbracelet/lipgloss"
)

type Request struct {
	Name   string
	Color  lipgloss.Color
	Active lipgloss.Color
}

var RequestCollection = []Request{
	{
		Name:  "GET",
		Color: styles.ColorGET,
	},
	{
		Name:  "POST",
		Color: styles.ColorPOST,
	},
	{
		Name:  "PUT",
		Color: styles.ColorPUT,
	},
	{
		Name:  "PATCH",
		Color: styles.ColorPATCH,
	},
	{
		Name:  "DELETE",
		Color: styles.ColorDELETE,
	},
	{
		Name:  "???",
		Color: styles.ColorUnknown,
	},
}
