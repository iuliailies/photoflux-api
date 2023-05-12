package config

import (
	"fmt"
	"time"

	cfg "github.com/Ozoniuss/configer"
)

func databaseOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "db-name", Shorthand: "", Value: "photoflux", ConfigKey: "database.name",
			Usage: "Specifies the name of the ports database."},
		{FlagName: "db-host", Shorthand: "", Value: "127.0.0.1", ConfigKey: "database.host",
			Usage: "Specifies the address on which the ports database listens for connections."},
		{FlagName: "db-port", Shorthand: "", Value: int32(5432), ConfigKey: "database.port",
			Usage: "Specifies the port on which the ports database listens for connections."},
		{FlagName: "db-user", Shorthand: "", Value: "photoflux", ConfigKey: "database.user",
			Usage: "Specifies the user which connects to the ports database."},
		{FlagName: "db-password", Shorthand: "", Value: "photoflux", ConfigKey: "database.password",
			Usage: "Specifies the password of the user which connects to the ports database."},
	}
}

func serverOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "server-address", Shorthand: "", Value: "127.0.0.1", ConfigKey: "server.address",
			Usage: "Specifies the address on which the ports service listens for incoming calls."},
		{FlagName: "server-port", Shorthand: "", Value: int32(8033), ConfigKey: "server.port",
			Usage: "Specifies the port on which the ports service listens for incoming calls."},
	}
}

func storageOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "storage-access-key", Shorthand: "", Value: "hello", ConfigKey: "storage.accesskey",
			Usage: "Access key of the admin minio client."},
		{FlagName: "storage-secret-key", Shorthand: "", Value: "myfriend", ConfigKey: "storage.secretkey",
			Usage: "Secret key of the admin minio client."},
		{FlagName: "stoarge-minio-address", Shorthand: "", Value: "127.0.0.1", ConfigKey: "storage.minioaddress",
			Usage: "Address of the minio server."},
		{FlagName: "storage-minio-port", Shorthand: "", Value: int32(9000), ConfigKey: "storage.minioport",
			Usage: "Port of the minio server."},
		{FlagName: "storage-user-policy-name", Shorthand: "", Value: "userpolicy", ConfigKey: "storage.userpolicyname",
			Usage: "Name of the policy assigned to all users."},
	}
}

func authOptions() []cfg.ConfigOption {
	return []cfg.ConfigOption{
		{FlagName: "auth-secret", Shorthand: "", Value: []byte("my_enconding_string"), ConfigKey: "auth.secret",
			Usage: "Secret for signing JWT toekns."},
		{FlagName: "access-token-lifetime", Shorthand: "", Value: 30 * time.Minute, ConfigKey: "auth.accesstokenlifetime",
			Usage: "Lifetime of the access token."},
		{FlagName: "minio-token-lifetime", Shorthand: "", Value: 30 * time.Minute, ConfigKey: "auth.miniotokenlifetime",
			Usage: "Lifetime of the minio token."},
	}
}

func allOptions() []cfg.ConfigOption {
	opts := make([]cfg.ConfigOption, 0)
	opts = append(opts, databaseOptions()...)
	opts = append(opts, serverOptions()...)
	opts = append(opts, storageOptions()...)
	opts = append(opts, authOptions()...)
	return opts
}

func ParseConfig() (Config, error) {
	c := newConfig()

	parserOptions := []cfg.ParserOption{
		cfg.WithConfigName("config"),
		cfg.WithConfigType("yml"),
		cfg.WithConfigPath("./configs"),
		cfg.WithEnvPrefix("PHOTOFLUX"),
		cfg.WithEnvKeyReplacer("_"),
		cfg.WithWriteFlag(),
	}

	err := cfg.NewConfig(&c, allOptions(), parserOptions...)
	if err != nil {
		return newConfig(), fmt.Errorf("could not create config: %w", err)
	}
	return c, nil
}
