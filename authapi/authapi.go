package authapi

import (
	"github.com/gorilla/mux"

  "github.com/matrix-org/dendrite/internal/config"
  "github.com/matrix-org/dendrite/authapi/routing"
)

// AddPublicRoutes sets up and registers HTTP handlers for the AuthAPI component.
func AddPublicRoutes(router *mux.Router, cfg *config.AuthAPI) {
  routing.Setup(router, cfg)
}
