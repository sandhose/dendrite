package config

import (
	"crypto/rsa"
)

const DefaultHashSecret = "DEFAULTSECRETPLEASECHANGE"

type AuthAPI struct {
	Matrix *Global `yaml:"-"`

	InternalAPI InternalAPIOptions `yaml:"internal_api"`

	Database DatabaseOptions `yaml:"database"`

	HashSecret       string          `yaml:"hash_secret"`
	OPPrivateKeyPath Path            `yaml:"op_private_key_path"`
	OPPrivateKey     *rsa.PrivateKey `yaml:"-"`
}

func (c *AuthAPI) Defaults() {
	c.InternalAPI.Listen = "http://localhost:7782"
	c.InternalAPI.Connect = "http://localhost:7782"
	c.Database.Defaults()
	c.Database.ConnectionString = "file:syncapi.db"
	c.HashSecret = DefaultHashSecret
	c.OPPrivateKeyPath = "op_key.pem"
}

func (c *AuthAPI) Verify(configErrs *ConfigErrors, isMonolith bool) {
	checkURL(configErrs, "auth_api.internal_api.listen", string(c.InternalAPI.Listen))
	checkURL(configErrs, "auth_api.internal_api.connect", string(c.InternalAPI.Connect))
	checkNotEmpty(configErrs, "auth_api.hash_secret", c.HashSecret)
	checkNotEmpty(configErrs, "auth_api.database", string(c.Database.ConnectionString))
}
