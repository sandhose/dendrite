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
  cfg *config.AuthAPI,
) {
  router.Handle("/token", httputil.MakeExternalAPI("oauth_token", func(req *http.Request) util.JSONResponse {
    return util.JSONResponse{
      Code: http.StatusNotImplemented,
      JSON: nil,
    }
  })).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

  router.Handle("/auth", httputil.MakeExternalAPI("oauth_auth", func(req *http.Request) util.JSONResponse {
    return util.JSONResponse{
      Code: http.StatusNotImplemented,
      JSON: nil,
    }
  })).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

  router.Handle("/userinfo", httputil.MakeExternalAPI("oidc_userinfo", func(req *http.Request) util.JSONResponse {
    return util.JSONResponse{
      Code: http.StatusNotImplemented,
      JSON: nil,
    }
  })).Methods(http.MethodGet, http.MethodOptions)
}
