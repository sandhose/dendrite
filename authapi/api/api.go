package api

import "context"

type Session struct {
	UserID   string `json:"user_id"`
	ClientID string `json:"client_id"`
}

type AuthInternalAPI interface {
	IntrospectAccessToken(ctx context.Context, req *AccessTokenIntrospectionRequest, res *AccessTokenIntrospectionResponse) error
}

type AccessTokenIntrospectionRequest struct {
	Token  string   `json:"token"`
	Scopes []string `json:"scopes"`
}

type AccessTokenIntrospectionResponse struct {
	*Session `json:"session"`
}
