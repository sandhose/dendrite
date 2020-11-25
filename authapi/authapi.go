package authapi

import (
	"github.com/gorilla/mux"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/token/jwt"
	"github.com/sirupsen/logrus"

	"github.com/matrix-org/dendrite/authapi/routing"
	"github.com/matrix-org/dendrite/authapi/storage"
	"github.com/matrix-org/dendrite/internal/config"
)

// AddPublicRoutes sets up and registers HTTP handlers for the AuthAPI component.
func AddPublicRoutes(router *mux.Router, cfg *config.AuthAPI) {
	db, err := storage.NewDatabase(&cfg.Database)
	if err != nil {
		logrus.WithError(err).Panicf("failed to create auth db")
	}

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

	provider := compose.Compose(
		config,
		storage,
		strategy,
		hasher,

		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2RefreshTokenGrantFactory,

		compose.OpenIDConnectExplicitFactory,
		compose.OpenIDConnectRefreshFactory,

		compose.OAuth2PKCEFactory,
	)

	routing.Setup(router, cfg, db, provider)
}
