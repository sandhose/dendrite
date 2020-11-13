package routing

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matrix-org/util"

	"github.com/matrix-org/dendrite/internal/config"
	"github.com/matrix-org/dendrite/internal/httputil"
)

func Setup(
	router *mux.Router,
	cfg *config.WellKnownAPI,
) {
	router.Handle("/openid-configuration", httputil.MakeExternalAPI("openid_configuration", func(req *http.Request) util.JSONResponse {
		return OpenIDConfiguration(req)
	})).Methods(http.MethodGet, http.MethodOptions)

	router.Handle("/jwks.json", httputil.MakeExternalAPI("jwks", func(req *http.Request) util.JSONResponse {
		return JWKS(req)
	})).Methods(http.MethodGet, http.MethodOptions)
}
