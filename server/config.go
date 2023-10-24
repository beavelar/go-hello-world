package server

type Config struct {
	Host     string
	Port     int
	Protocol string
}

var GrpcServerConfig *Config
var HttpServerConfig *Config
var WebsocketServerConfig *Config
