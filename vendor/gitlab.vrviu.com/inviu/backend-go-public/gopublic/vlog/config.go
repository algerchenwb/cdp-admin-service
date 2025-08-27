package vlog

import (
	"gitlab.vrviu.com/inviu/backend-go-public/gopublic/vlog/tracing"
)

type Config tracing.Config

func NewDevelopmentConfig(appname string, filename string) Config {
	return Config(tracing.NewDevelopmentConfig(appname, filename))
}

func NewProductionConfig(appname string, filename string) Config {
	return Config(tracing.NewProductionConfig(appname, filename))
}

func NewStdConfig() Config {
	return Config(tracing.NewStdConfig())
}

func NewConfig(filename string) (Config, error) {
	config, err := tracing.NewConfig(filename)
	return Config(config), err
}

type Option func(*tracing.Config)

func WithGlobalCallerSkip(n int) Option {
	return func(c *tracing.Config) {
		c.GlobalCallerSkip = n
	}
}
