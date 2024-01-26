package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env                   string
	AgeServiceUrl         string
	NationalityServiceUrl string
	GenderServiceUrl      string
	CtxTimeout            int
	Storage               StorageConfig
	HTTPServer            HTTPServer
}

type HTTPServer struct {
	Address     string
	Timeout     time.Duration
	IdleTimeout time.Duration
}

type StorageConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func MustLoad() *Config {
	var cfg Config
	var err error

	cfg.Env = os.Getenv("PS_ENV")

	cfg.AgeServiceUrl = os.Getenv("PS_AGE_URL")
	cfg.GenderServiceUrl = os.Getenv("PS_GENDER_URL")
	cfg.NationalityServiceUrl = os.Getenv("PS_NATIONALITY_URL")

	cfg.CtxTimeout, err = strconv.Atoi(loadConfig("PS_CTX_TIMEOUT"))
	if err != nil {
		panic(fmt.Sprintf("cannot load ctx timeout config: %s", err))
	}

	cfg.Storage.Host = loadConfig("PS_PG_DB_HOST")
	cfg.Storage.Port, err = strconv.Atoi(loadConfig("PS_PG_DB_PORT"))
	if err != nil {
		panic(fmt.Sprintf("cannot load db port config: %s", err))
	}
	cfg.Storage.DBName = loadConfig("PS_PG_DB_NAME")
	cfg.Storage.User = loadConfig("PS_PG_DB_USER")
	cfg.Storage.Password = loadConfig("PS_PG_DB_PASS")

	cfg.HTTPServer.Address = loadConfig("PS_HTTP_SERVER")
	cfg.HTTPServer.Timeout, err = time.ParseDuration(loadConfig("PS_HTTP_TIMEOUT"))
	if err != nil {
		panic(fmt.Sprintf("cannot load server timeout config: %s", err))
	}
	cfg.HTTPServer.IdleTimeout, err = time.ParseDuration(loadConfig("PS_HTTP_IDLE_TIMEOUT"))
	if err != nil {
		panic(fmt.Sprintf("cannot load server idle timeout config: %s", err))
	}

	return &cfg
}

func loadConfig(name string) string {
	cfg, exists := os.LookupEnv(name)

	if !exists {
		panic(fmt.Sprintf("env variable does not exist: %s", name))
	}

	return cfg
}
