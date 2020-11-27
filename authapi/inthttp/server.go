package inthttp

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matrix-org/dendrite/authapi/api"
	"github.com/matrix-org/dendrite/internal/httputil"
	"github.com/matrix-org/util"
)

func AddRoutes(internalAPIMux *mux.Router, s api.AuthInternalAPI) {
	internalAPIMux.Handle(IntrospectAccessTokenPath,
		httputil.MakeInternalAPI("introspectAccessToken", func(req *http.Request) util.JSONResponse {
			request := api.AccessTokenIntrospectionRequest{}
			response := api.AccessTokenIntrospectionResponse{}
			if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
				return util.MessageResponse(http.StatusBadRequest, err.Error())
			}
			if err := s.IntrospectAccessToken(req.Context(), &request, &response); err != nil {
				return util.ErrorResponse(err)
			}
			return util.JSONResponse{Code: http.StatusOK, JSON: &response}
		}),
	)
}
