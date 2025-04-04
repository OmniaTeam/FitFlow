package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

var (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

type Config struct {
	Env            string         `yaml:"env" env-required:"true"`
	HTTPServer     HTTPServer     `yaml:"http_server"`
	PostgresConfig PostgresConfig `yaml:"postgres"`
	JwtConfig      JwtConfig      `yaml:"jwt"`
}

type HTTPServer struct {
	Address         string        `yaml:"address" env-default:"localhost:8080"`
	Timeout         time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env-default:"60s"`
	GinModeRelease  bool          `yaml:"gin_mode" env-default:"false"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env-default:"5s"`
	Oauth           Oauth         `yaml:"oauth"`
	Routes          struct {
		PublicRoutes  []RouteConfig `yaml:"public"`
		PrivateRoutes []RouteConfig `yaml:"private"`
	} `yaml:"routes"`
}

type Oauth struct {
	Google   OauthCredentials `yaml:"google"`
	Vkid     OauthCredentials `yaml:"vkid"`
	Redirect string           `yaml:"redirect"`
}

type OauthCredentials struct {
	ClientID string `yaml:"client_id"`
	SecretID string `yaml:"secret_id"`
}

type RouteConfig struct {
	Path   string   `yaml:"path"`
	Method string   `yaml:"method"`
	Roles  []string `yaml:"roles"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

type JwtConfig struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env-required:"true"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-required:"true"`
	Secret          string        `yaml:"secret" env-required:"true"`
}

func MustLoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	return &config
}
