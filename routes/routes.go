// routes/routes.go
package routes

import (
	"net/http"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Enregistrer toutes les routes
	RegisterAuthRoutes(mux)
	RegisterAPIRoutes(mux)
	RegisterAdminRoutes(mux)

	return mux
}
