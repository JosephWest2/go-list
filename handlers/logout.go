
package handlers

import (
	"net/http"

	"josephwest2.com/go-list/components"
)

func LogoutPageHandler(w http.ResponseWriter, r *http.Request) {
	RenderPage(components.LogoutPage(), w)
}
