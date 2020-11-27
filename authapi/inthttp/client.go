package inthttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/matrix-org/dendrite/internal/httputil"
	"github.com/matrix-org/dendrite/authapi/api"
	"github.com/opentracing/opentracing-go"
)

const (
	IntrospectAccessTokenPath = "/authapi/introspectAccessToken"
)

func NewAuthAPIClient(
	apiURL string,
	httpClient *http.Client,
) (api.AuthInternalAPI, error) {
	if httpClient == nil {
		return nil, errors.New("NewUserAPIClient: httpClient is <nil>")
	}
	return &httpAuthInternalAPI{
		apiURL:     apiURL,
		httpClient: httpClient,
	}, nil
}

type httpAuthInternalAPI struct {
	apiURL     string
	httpClient *http.Client
}

func (h *httpAuthInternalAPI) IntrospectAccessToken(ctx context.Context, req *api.AccessTokenIntrospectionRequest, res *api.AccessTokenIntrospectionResponse) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "IntrospectAccessToken")
	defer span.Finish()

	apiURL := h.apiURL + IntrospectAccessTokenPath
	return httputil.PostJSON(ctx, span, h.httpClient, apiURL, req, res)
}
