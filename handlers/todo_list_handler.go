package handlers

import (
	"net/http"

	"josephwest2.com/go-list/components"
)

func TodoListPageHandler(w http.ResponseWriter, r *http.Request) {
	RenderPage(components.TodoListPage(), w)
}
