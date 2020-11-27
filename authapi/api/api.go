package api

import "context"

type AuthInternalAPI interface {
	IntrospectAccessToken(ctx context.Context, req *AccessTokenIntrospectionRequest, res *AccessTokenIntrospectionResponse) error
}

type AccessTokenIntrospectionRequest struct {
	Token  string   `json:"token"`
	Scopes []string `json:"scopes"`
}

type AccessTokenIntrospectionResponse struct {
	Session *Session `json:"session"`
	Client  *Client  `json:"client"`
}
