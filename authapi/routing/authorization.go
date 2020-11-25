package routing

import (
	"net/http"

	"github.com/ory/fosite"
	"github.com/ory/fosite/handler/openid"
	log "github.com/sirupsen/logrus"
)

func Authorize(rw http.ResponseWriter, req *http.Request, oauth2Provider fosite.OAuth2Provider) {
	ctx := req.Context()

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

	// TODO: implement the fosite.Session & fosite/handler/openid.Session interfaces
	session := openid.NewDefaultSession()
	session.Subject = "bob"
	session.Claims.Subject = "bob"

	response, err := oauth2Provider.NewAuthorizeResponse(ctx, ar, session)
	if err != nil {
		log.WithError(err).WithContext(ctx).Error("Could not fullfil authorization request")
		oauth2Provider.WriteAuthorizeError(rw, ar, err)
		return
	}

	oauth2Provider.WriteAuthorizeResponse(rw, ar, response)
}
