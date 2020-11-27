package api

import (
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/token/jwt"
)

type Session struct {
	LocalPart string                         `json:"localpart"`
	DeviceID  string                         `json:"device_id"`
	ExpiresAt map[fosite.TokenType]time.Time `json:"expires_at"`
}

func NewSession(userID string, deviceID string) *Session {
	return &Session{
		LocalPart: userID,
		DeviceID:  deviceID,
		ExpiresAt: make(map[fosite.TokenType]time.Time),
	}
}

func NewEmptySession() *Session {
	return NewSession("", "")
}

func (s *Session) SetExpiresAt(key fosite.TokenType, exp time.Time) {
	s.ExpiresAt[key] = exp
}

func (s *Session) GetExpiresAt(key fosite.TokenType) time.Time {
	if _, ok := s.ExpiresAt[key]; !ok {
		return time.Time{}
	}

	return s.ExpiresAt[key]
}

func (s *Session) GetUsername() string {
	return "" // TODO
}

func (s *Session) GetSubject() string {
	return s.LocalPart
}

func (s *Session) Clone() fosite.Session {
	if s == nil {
		return nil
	}

	t := &Session{
		LocalPart: s.LocalPart,
		DeviceID:  s.DeviceID,
		ExpiresAt: make(map[fosite.TokenType]time.Time),
	}

	for key, value := range s.ExpiresAt {
		t.ExpiresAt[key] = value
	}

	return t
}

func (s *Session) IDTokenHeaders() *jwt.Headers {
	return &jwt.Headers{}
}

func (s *Session) IDTokenClaims() *jwt.IDTokenClaims {
	return &jwt.IDTokenClaims{
		Subject: s.GetSubject(),
	}
}
