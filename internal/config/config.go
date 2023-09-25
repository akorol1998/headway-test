package config

import (
	"os"
)

const (
	dbHost           = "DB_HOST"
	dbPort           = "DB_PORT"
	dbUser           = "DB_USER"
	dbName           = "DB_NAME"
	dbPass           = "DB_PASS"
	serviceHost      = "SERVICE_HOST"
	servicePort      = "SERVICE_PORT"
	logLevel         = "LOG_LEVEL"
	envName          = "ENVIRONMENT"
	providerFilePath = "PROVIDER_FILE_PATH"
	storesFilePath   = "STORES_FILE_PATH"
)

type ConfigDB struct {
	Host     string
	Port     string
	User     string
	DbName   string
	Password string
}

type ConfigService struct {
	Host string
	Port string
}

type Config struct {
	Database         ConfigDB
	Service          ConfigService
	LogLevel         string
	Environment      string
	ProviderFilePath string
	StoresFilePath   string
}

// Load loads env variables
func Load() *Config {
	return &Config{
		Database:         database(),
		Service:          service(),
		LogLevel:         logger(),
		Environment:      environment(),
		ProviderFilePath: provider(),
		StoresFilePath:   stores(),
	}
}

func database() ConfigDB {
	conf := ConfigDB{}
	conf.Host = os.Getenv(dbHost)
	if len(conf.Host) == 0 {
		conf.Host = "localhost"
	}
	conf.Port = os.Getenv(dbPort)
	if len(conf.Port) == 0 {
		conf.Port = "5432"
	}
	conf.User = os.Getenv(dbUser)
	if len(conf.User) == 0 {
		conf.User = "postgres"
	}
	conf.DbName = os.Getenv(dbName)
	if len(conf.DbName) == 0 {
		conf.DbName = "postgres"
	}
	conf.Password = os.Getenv(dbPass)
	if len(conf.Password) == 0 {
		conf.Password = "adminlol"
	}
	return conf
}

func service() ConfigService {
	conf := ConfigService{}
	conf.Host = os.Getenv(serviceHost)
	if len(conf.Host) == 0 {
		conf.Host = "localhost"
	}
	conf.Port = os.Getenv(servicePort)
	if len(conf.Port) == 0 {
		conf.Port = "8080"
	}
	return conf
}

func logger() string {
	ll := os.Getenv(logLevel)
	if len(ll) == 0 {
		ll = "info"
	}
	return ll
}

func environment() string {
	env := os.Getenv(envName)
	if len(env) == 0 {
		env = "development"
	}
	return env
}

func provider() string {
	env := os.Getenv(providerFilePath)
	if len(env) == 0 {
		env = "./assets/providers.json"
	}
	return env
}

func stores() string {
	env := os.Getenv(storesFilePath)
	if len(env) == 0 {
		env = "./assets/stores.json"
	}
	return env
}
