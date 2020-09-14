package conf

import (
	"sync"

	"github.com/isayme/go-config"
	duration "github.com/isayme/go-duration"
	"github.com/isayme/go-logger"
)

type Logger struct {
	Level string `json:"level" yaml:"level"`
}

type Server struct {
	Addr string `json:"addr" yaml:"addr"`
}

type Local struct {
	WebsocketAddr string `json:"ws_addr" yaml:"ws_addr"`
}

type ServiceConfig struct {
	Name    string            `json:"name" yaml:"name"`
	Timeout duration.Duration `json:"timeout" yaml:"timeout"`

	// for server
	RemoteAddress string `json:"remote_addr" yaml:"remote_addr"`

	// for local
	LocalAddress string `json:"local_addr" yaml:"local_addr"`
}

type Config struct {
	Logger Logger `json:"logger" yaml:"logger"`

	Local  Local  `json:"local" yaml:"local"`
	Server Server `json:"server" yaml:"server"`

	Services []ServiceConfig `json:"services" yaml:"services"`
}

var once sync.Once
var globalConfig Config

func Get() *Config {
	config.Parse(&globalConfig)
	once.Do(func() {
		logger.SetLevel(globalConfig.Logger.Level)
	})
	return &globalConfig
}
