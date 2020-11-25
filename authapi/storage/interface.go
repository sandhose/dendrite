package storage

import (
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/storage"
	"golang.org/x/net/context"

	"github.com/matrix-org/dendrite/authapi/api"
)

type Database interface {
	GetClient(ctx context.Context, id string) (*api.Client, error)
	CreateClient(ctx context.Context, client *api.Client) error
	UpdateClient(ctx context.Context, client *api.Client) error
}

type WrappedStorage struct {
	Database
	*storage.MemoryStore
}

func WrapDatabase(db Database) *WrappedStorage {
	return &WrappedStorage{
		Database:    db,
		MemoryStore: storage.NewMemoryStore(),
	}
}

// Implement fosite.Storage

func (s *WrappedStorage) GetClient(ctx context.Context, id string) (fosite.Client, error) {
	return s.Database.GetClient(ctx, id)
}

func (s *WrappedStorage) ClientAssertionJWTValid(ctx context.Context, jti string) error {
	return s.MemoryStore.ClientAssertionJWTValid(ctx, jti)
}

func (s *WrappedStorage) SetClientAssertionJWT(ctx context.Context, jti string, exp time.Time) error {
	return s.MemoryStore.SetClientAssertionJWT(ctx, jti, exp)
}

// Implement oauth2.AuthorizeCodeStorage

func (s *WrappedStorage) CreateAuthorizeCodeSession(ctx context.Context, code string, request fosite.Requester) error {
	return s.MemoryStore.CreateAuthorizeCodeSession(ctx, code, request)
}

func (s *WrappedStorage) GetAuthorizeCodeSession(ctx context.Context, code string, session fosite.Session) (fosite.Requester, error) {
	return s.MemoryStore.GetAuthorizeCodeSession(ctx, code, session)
}

func (s *WrappedStorage) InvalidateAuthorizeCodeSession(ctx context.Context, code string) error {
	return s.MemoryStore.InvalidateAuthorizeCodeSession(ctx, code)
}

// Implement oauth2.AccessTokenStorage

func (s *WrappedStorage) CreateAccessTokenSession(ctx context.Context, signature string, request fosite.Requester) error {
	return s.MemoryStore.CreateAccessTokenSession(ctx, signature, request)
}

func (s *WrappedStorage) GetAccessTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.MemoryStore.GetAccessTokenSession(ctx, signature, session)
}

func (s *WrappedStorage) DeleteAccessTokenSession(ctx context.Context, signature string) error {
	return s.MemoryStore.DeleteAccessTokenSession(ctx, signature)
}

// Implement oauth2.RefreshTokenStorage

func (s *WrappedStorage) CreateRefreshTokenSession(ctx context.Context, signature string, request fosite.Requester) error {
	return s.MemoryStore.CreateRefreshTokenSession(ctx, signature, request)
}

func (s *WrappedStorage) GetRefreshTokenSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.MemoryStore.GetRefreshTokenSession(ctx, signature, session)
}

func (s *WrappedStorage) DeleteRefreshTokenSession(ctx context.Context, signature string) error {
	return s.MemoryStore.DeleteRefreshTokenSession(ctx, signature)
}

// Implement oauth2.TokenRevocationStorage

func (s *WrappedStorage) RevokeRefreshToken(ctx context.Context, requestID string) error {
	return s.MemoryStore.RevokeRefreshToken(ctx, requestID)
}

func (s *WrappedStorage) RevokeAccessToken(ctx context.Context, requestID string) error {
	return s.MemoryStore.RevokeAccessToken(ctx, requestID)
}

// Implement openid.OpenIDConnectRequestStorage

func (s *WrappedStorage) CreateOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) error {
	return s.MemoryStore.CreateOpenIDConnectSession(ctx, authorizeCode, requester)
}

func (s *WrappedStorage) GetOpenIDConnectSession(ctx context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	return s.MemoryStore.GetOpenIDConnectSession(ctx, authorizeCode, requester)
}

func (s *WrappedStorage) DeleteOpenIDConnectSession(ctx context.Context, authorizeCode string) error {
	return s.MemoryStore.DeleteOpenIDConnectSession(ctx, authorizeCode)
}

// Implement pkce.PKCERequestStorage

func (s *WrappedStorage) GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error) {
	return s.MemoryStore.GetPKCERequestSession(ctx, signature, session)
}

func (s *WrappedStorage) CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error {
	return s.MemoryStore.CreatePKCERequestSession(ctx, signature, requester)
}

func (s *WrappedStorage) DeletePKCERequestSession(ctx context.Context, signature string) error {
	return s.MemoryStore.DeletePKCERequestSession(ctx, signature)
}
