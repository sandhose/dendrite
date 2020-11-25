package routing

import (
	"net/http"

	"github.com/matrix-org/util"
)

type jwk struct {
	Use string `json:"use,omitempty"`
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

func JWKS(req *http.Request) util.JSONResponse {
	return util.JSONResponse{
		Code: http.StatusOK,
		JSON: jwks{
			Keys: []jwk{},
		},
	}
}
