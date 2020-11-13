package config

type AuthAPI struct {
	Matrix *Global `yaml:"-"`

	InternalAPI InternalAPIOptions `yaml:"internal_api"`
}

func (c *AuthAPI) Defaults() {
	c.InternalAPI.Listen = "http://localhost:7782"
	c.InternalAPI.Connect = "http://localhost:7782"
}

func (c *AuthAPI) Verify(configErrs *ConfigErrors, isMonolith bool) {
	checkURL(configErrs, "auth_api.internal_api.listen", string(c.InternalAPI.Listen))
	checkURL(configErrs, "auth_api.internal_api.connect", string(c.InternalAPI.Connect))
}
