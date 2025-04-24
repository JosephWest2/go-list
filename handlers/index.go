package handlers

import (
	"net/http"

	"josephwest2.com/go-list/components"
)

func IndexPageHandler(w http.ResponseWriter, r *http.Request) {
	RenderPage(components.HomePage(), w)
}
