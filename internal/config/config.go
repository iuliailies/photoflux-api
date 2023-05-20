package config

import "time"

type Config struct {
	Server        Server
	Database      Database
	MongoDatabase MongoDatabase
	ApiPaths      ApiPaths
	Storage       Storage
	Auth          Auth
	Notifications Notifications
}

type Database struct {
	Host     string
	Port     int32
	User     string
	Name     string
	Password string
}

type MongoDatabase struct {
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

type Auth struct {
	Secret              []byte
	AccessTokenLifetime time.Duration
	MinioTokenLifetime  time.Duration
}

type ApiPaths struct {
	Photos     string
	Users      string
	Categories string
}

type Storage struct {
	AccessKey      string
	SecretKey      string
	MinioAddress   string
	MinioPort      int32
	UserPolicyName string
}

type Notifications struct {
	RabbitMQ RabbitMQ
}

type RabbitMQ struct {
	User     string
	Password string
	Host     string
	Port     int32
	Vhost    string
	Exchange string
	Queue    string
}

func newConfig() Config {
	c := Config{}
	return c
}
