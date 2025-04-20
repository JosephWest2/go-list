package handlers

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
	"josephwest2.com/go-list/components"
)

func RenderPage(component templ.Component, w http.ResponseWriter) {
	components.Layout(component).Render(context.Background(), w)
}
