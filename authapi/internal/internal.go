package internal

import (
	"context"

	"github.com/ory/fosite"

	"github.com/matrix-org/dendrite/authapi/api"
)

type AuthInternalAPI struct {
	fosite.OAuth2Provider
}

func (a *AuthInternalAPI) IntrospectAccessToken(ctx context.Context, req *api.AccessTokenIntrospectionRequest, res *api.AccessTokenIntrospectionResponse) error {
	fositeSession := &fosite.DefaultSession{}
	_, ar, err := a.IntrospectToken(ctx, req.Token, fosite.AccessToken, fositeSession, req.Scopes...)
	if err != nil {
		return err
	}

	res.Session = &api.Session{
		ClientID: ar.GetClient().GetID(),
		UserID:   fositeSession.GetSubject(),
	}

	return nil
}
