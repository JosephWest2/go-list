package handlers

import (
	"net/http"

	"josephwest2.com/go-list/components"
)

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	RenderPage(components.RegisterPage(), w)
}
