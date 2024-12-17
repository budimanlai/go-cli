package gocli

type Config struct {
	configName string
	configType ConfigType
}

func NewConfig() Config {
	cfg := Config{}

	return cfg
}
func (c *Config) SetConfigName(name string) {
	c.configName = name
}

func (c *Config) SetConfigType(name ConfigType) {
	c.configType = name
}
