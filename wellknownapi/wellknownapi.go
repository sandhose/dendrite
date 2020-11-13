package wellknownapi

import (
	"github.com/gorilla/mux"

	"github.com/matrix-org/dendrite/internal/config"
	"github.com/matrix-org/dendrite/wellknownapi/routing"
)

// AddPublicRoutes sets up and registers HTTP handlers for the WellKnownAPI component.
func AddPublicRoutes(router *mux.Router, cfg *config.WellKnownAPI) {
	routing.Setup(router, cfg)
}
