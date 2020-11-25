package routing

import (
	"net/http"

	"github.com/ory/fosite"
)

func Token(rw http.ResponseWriter, req *http.Request, oauth2Provider fosite.OAuth2Provider) {
	ctx := req.Context()

	session := &fosite.DefaultSession{}

	accessRequest, err := oauth2Provider.NewAccessRequest(ctx, req, session)
	if err != nil {
		oauth2Provider.WriteAccessError(rw, accessRequest, err)
		return
	}

	response, err := oauth2Provider.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		oauth2Provider.WriteAccessError(rw, accessRequest, err)
		return
	}

	oauth2Provider.WriteAccessResponse(rw, accessRequest, response)
}
