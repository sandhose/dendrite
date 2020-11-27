package routing

import (
	"net/http"

	"github.com/ory/fosite"

	"github.com/matrix-org/dendrite/authapi/api"
	"github.com/matrix-org/dendrite/authapi/storage"
	"github.com/matrix-org/util"
)

func Authorize(rw http.ResponseWriter, req *http.Request, oauth2Provider fosite.OAuth2Provider, db storage.Database) {
	ctx := req.Context()
	log := util.GetLogger(ctx)

	ar, err := oauth2Provider.NewAuthorizeRequest(ctx, req)
	if err != nil {
		log.WithError(err).WithContext(ctx).Error("Invalid authorization request")
		oauth2Provider.WriteAuthorizeError(rw, ar, err)
		return
	}

	for _, scope := range ar.GetRequestedScopes() {
		ar.GrantScope(scope)
	}

	// TODO: actually authorize users
	session := api.NewSession("foo", "bar")

	response, err := oauth2Provider.NewAuthorizeResponse(ctx, ar, session)
	if err != nil {
		log.WithError(err).WithContext(ctx).Error("Could not fullfil authorization request")
		oauth2Provider.WriteAuthorizeError(rw, ar, err)
		return
	}

	if err := db.CreateSession(ctx, ar.GetID(), session); err != nil {
		// TODO: write error
		return
	}

	oauth2Provider.WriteAuthorizeResponse(rw, ar, response)
}
