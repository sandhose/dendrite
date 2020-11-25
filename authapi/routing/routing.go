package routing

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matrix-org/util"
	"github.com/ory/fosite"

	"github.com/matrix-org/dendrite/authapi/storage"
	"github.com/matrix-org/dendrite/internal/config"
	"github.com/matrix-org/dendrite/internal/httputil"
)

func Setup(
	router *mux.Router,
	cfg *config.AuthAPI,
	database storage.Database,
	oauth2Provider fosite.OAuth2Provider,
) {
	router.Handle("/token", httputil.MakeHTMLAPI("oauth_token", func(rw http.ResponseWriter, req *http.Request) *util.JSONResponse {
		Token(rw, req, oauth2Provider)
		return nil
	})).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	router.Handle("/auth", httputil.MakeHTMLAPI("oauth_auth", func(rw http.ResponseWriter, req *http.Request) *util.JSONResponse {
		Authorize(rw, req, oauth2Provider)
		return nil
	})).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	router.Handle("/userinfo", httputil.MakeExternalAPI("oidc_userinfo", func(req *http.Request) util.JSONResponse {
		return util.JSONResponse{
			Code: http.StatusNotImplemented,
			JSON: nil,
		}
	})).Methods(http.MethodGet, http.MethodOptions)

	router.Handle("/clients/register", httputil.MakeExternalAPI("oauth_client_register", func(req *http.Request) util.JSONResponse {
		return RegisterClient(req, database)
	})).Methods(http.MethodPost, http.MethodOptions)
}
