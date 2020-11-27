package authapi

import (
	"github.com/gorilla/mux"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/token/jwt"
	"github.com/sirupsen/logrus"

	"github.com/matrix-org/dendrite/authapi/api"
	"github.com/matrix-org/dendrite/authapi/internal"
	"github.com/matrix-org/dendrite/authapi/inthttp"
	"github.com/matrix-org/dendrite/authapi/routing"
	"github.com/matrix-org/dendrite/authapi/storage"
	"github.com/matrix-org/dendrite/internal/config"
)

func NewProvider(db storage.Database, cfg *config.AuthAPI) fosite.OAuth2Provider {
	config := &compose.Config{
		EnforcePKCEForPublicClients: true,
		RefreshTokenScopes:          []string{}, // Generate a refresh token for all scopes
		IDTokenIssuer:               cfg.Matrix.Issuer,
	}

	secret := []byte(cfg.HashSecret)

	strategy := &compose.CommonStrategy{
		CoreStrategy: compose.NewOAuth2HMACStrategy(config, secret, nil),
		JWTStrategy: &jwt.RS256JWTStrategy{
			PrivateKey: cfg.OPPrivateKey,
		},
		OpenIDConnectTokenStrategy: compose.NewOpenIDConnectStrategy(config, cfg.OPPrivateKey),
	}

	hasher := &fosite.BCrypt{}

	storage := storage.WrapDatabase(db)

	return compose.Compose(
		config,
		storage,
		strategy,
		hasher,

		compose.OAuth2TokenIntrospectionFactory,

		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2RefreshTokenGrantFactory,

		compose.OpenIDConnectExplicitFactory,
		compose.OpenIDConnectRefreshFactory,

		compose.OAuth2PKCEFactory,
	)
}

func Init(cfg *config.AuthAPI) (storage.Database, fosite.OAuth2Provider) {
	db, err := storage.NewDatabase(&cfg.Database)
	if err != nil {
		logrus.WithError(err).Panicf("failed to create auth db")
	}

	provider := NewProvider(db, cfg)

	return db, provider
}

// AddInternalRoutes registers HTTP handlers for the internal API. Invokes functions
// on the given input API.
func AddInternalRoutes(router *mux.Router, intAPI api.AuthInternalAPI) {
	inthttp.AddRoutes(router, intAPI)
}

// AddPublicRoutes sets up and registers HTTP handlers for the AuthAPI component.
func AddPublicRoutes(router *mux.Router, cfg *config.AuthAPI, db storage.Database, provider fosite.OAuth2Provider) {
	routing.Setup(router, cfg, db, provider)
}

func NewInternalAPI(_cfg *config.AuthAPI, db storage.Database, provider fosite.OAuth2Provider) api.AuthInternalAPI {
	return &internal.AuthInternalAPI{
		OAuth2Provider: provider,
		Database:       db,
	}
}
