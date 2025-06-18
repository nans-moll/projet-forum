// routes/routes.go
package routes

import (
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	RegisterAuthRoutes(mux)
	RegisterAPIRoutes(mux)
	RegisterAdminRoutes(mux)

	return mux
}
