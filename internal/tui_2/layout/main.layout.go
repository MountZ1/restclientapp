package layout

import (
	"restclient/internal/collections"
	"restclient/internal/tui_2/components"
)

func MainLayout() string {
	request := components.List(collections.RequestArray[:])

	return request.View()
}
