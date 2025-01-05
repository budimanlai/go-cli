package gocli

type CliOptions struct {
	ConfigFile     []string
	RuntimePath    string
	AutoLoadConfig bool
	AppName        string
	Version        string
}
