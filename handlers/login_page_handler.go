
package handlers

import (
	"net/http"

	"josephwest2.com/go-list/components"
)

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	RenderPage(components.LoginPage(), w)
}
