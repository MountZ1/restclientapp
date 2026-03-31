package layout

import (
	"restclient/internal/collections"
	"restclient/internal/tui/components"
)

func MainLayout() string {
	request := components.List(collections.RequestArray[:])

	return request.View()
}
