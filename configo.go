package configo

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type App struct {
	Name    string `yaml:"name" env-required:"true"`
	Version string `yaml:"version" env-required:"true"`
}

type Logger struct {
	Level         int    `yaml:"level" env-default:"0"`
	Dir           string `yaml:"dir" env-default:"logs"`
	MaxSize       int    `yaml:"maxSize" env-default:"10"`
	MaxBackups    int    `yaml:"maxBackups" env-default:"3"`
	MaxAge        int    `yaml:"maxAge" env-default:"365"`
	Compress      bool   `yaml:"compress" env-default:"true"`
	RotationTime  string `yaml:"rotationTime" env-default:"24h"`
	ConsoleLevel  int    `yaml:"consoleLevel" env-default:"0"`
	FileLevel     int    `yaml:"fileLevel" env-default:"0"`
	EnableConsole bool   `yaml:"enableConsole" env-default:"true"`
	EnableFile    bool   `yaml:"enableFile" env-default:"true"`
	TimeFormat    string `yaml:"timeFormat" env-default:"2006-01-02T15:04:05.000Z07:00"`
}

type Database struct {
	Type          string        `yaml:"type" env-required:"true"`
	Host          string        `yaml:"host" env-required:"true"`
	Port          int           `yaml:"port" env-required:"true"`
	Name          string        `yaml:"name" env-required:"true"`
	User          string        `yaml:"user" env-required:"true"`
	Password      string        `yaml:"password" env-required:"true"`
	Schema        string        `yaml:"schema" env-default:"public"`
	MigrationPath string        `yaml:"migrationPath" env-required:"true"`
	MaxAttempts   int           `yaml:"maxAttempts" env-required:"true"`
	AttemptDelay  time.Duration `yaml:"attemptDelay" env-required:"true"`
}

type Redis struct {
	Host string `yaml:"host" env-required:"true"`
	Port int    `yaml:"port" env-required:"true"`
	Db   int    `yaml:"db" `
}

type Sentry struct {
	Host string `yaml:"host" env-required:"true"`
	Key  string `yaml:"key" env-required:"true"`
}

type Service struct {
	Port uint16 `yaml:"port" env-required:"true"`
}

type S3 struct {
	Endpoint  string `yaml:"endpoint" env-required:"true"`
	Region    string `yaml:"region" env-required:"true"`
	AccessKey string `yaml:"accessKey" env-required:"true"`
	SecretKey string `yaml:"secretKey" env-required:"true"`
}

type Ws struct {
	Port               int     `yaml:"port" env-required:"true"`
	MaxOneIpConnection int     `yaml:"maxOneIpConnection" env-required:"true"`
	Session            session `yaml:"session"`
}

type session struct {
	MinPingDuration       time.Duration `yaml:"minPingDuration" env-required:"true"`
	MaxPingDuration       time.Duration `yaml:"maxPingDuration" env-required:"true"`
	MaxInactivityDuration time.Duration `yaml:"maxInactivityDuration" env-required:"true"`
}

func MustLoad[TConfig any]() *TConfig {
	path := fetchConfigPath()

	if path == "" {
		panic("Путь конфига не найден")
	}

	if _, err := os.Stat(path); err != nil {
		panic("Файл конфига не найден")
	}

	var cfg TConfig

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Ошибка загрузки конфига: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var result string
	flag.StringVar(&result, "config", "", "Путь до файла конфига")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return result
}
