package constant

import (
	"image/color"
	"restclient/internal/tui/styles"
)

type Request struct {
	Name   string
	Color  color.Color
	Active color.Color
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
