package internal

import (
	"context"

	"github.com/ory/fosite"

	"github.com/matrix-org/dendrite/authapi/api"
	"github.com/matrix-org/dendrite/authapi/storage"
)

type AuthInternalAPI struct {
	fosite.OAuth2Provider
	storage.Database
}

func (a *AuthInternalAPI) IntrospectAccessToken(ctx context.Context, req *api.AccessTokenIntrospectionRequest, res *api.AccessTokenIntrospectionResponse) error {
	_, ar, err := a.IntrospectToken(ctx, req.Token, fosite.AccessToken, nil, req.Scopes...)
	if err != nil {
		return err
	}

	res.Client = ar.GetClient().(*api.Client)
	if res.Session, err = a.GetSession(ctx, ar.GetID()); err != nil {
		return err
	}

	return nil
}
