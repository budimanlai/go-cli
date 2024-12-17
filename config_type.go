package gocli

type ConfigType struct {
	slug string
}

func (c ConfigType) String() string {
	return c.slug
}

var (
	Json ConfigType = ConfigType{"json"}
	Yaml ConfigType = ConfigType{"yaml"}
)
