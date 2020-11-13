package config

type WellKnownAPI struct {
	Matrix *Global `yaml:"-"`

	InternalAPI InternalAPIOptions `yaml:"internal_api"`
}

func (c *WellKnownAPI) Defaults() {
	c.InternalAPI.Listen = "http://localhost:7783"
	c.InternalAPI.Connect = "http://localhost:7783"
}

func (c *WellKnownAPI) Verify(configErrs *ConfigErrors, isMonolith bool) {
	checkURL(configErrs, "well_known_api.internal_api.listen", string(c.InternalAPI.Listen))
	checkURL(configErrs, "well_known_api.internal_api.connect", string(c.InternalAPI.Connect))
}
