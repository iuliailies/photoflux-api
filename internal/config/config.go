package config

type Config struct {
	Server   Server
	Database Database
	ApiPaths ApiPaths
}

type Database struct {
	Host     string
	Port     int32
	User     string
	Name     string
	Password string
}

type Server struct {
	Address string
	Port    int32
}

type ApiPaths struct {
	Photos     string
	Users      string
	Categories string
}

func newConfig() Config {
	c := Config{}
	return c
}
