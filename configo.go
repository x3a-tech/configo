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

type KafkaProducer struct {
	Brokers      []string      `yaml:"brokers" env-separator:"," env-required:"true"`
	RequiredAcks int           `yaml:"requiredAcks" env-default:"1"` // Уровень подтверждения: 0=None, 1=Leader, -1=All
	Async        bool          `yaml:"async" env-default:"false"`
	BatchSize    int           `yaml:"batchSize" env-default:"100"`
	BatchTimeout time.Duration `yaml:"batchTimeout" env-default:"1s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"10s"`
	MaxAttempts  int           `yaml:"maxAttempts" env-default:"3"`
}

type KafkaConsumer struct {
	Brokers     []string `yaml:"brokers" env-separator:"," env-required:"true"`
	GroupID     string   `yaml:"groupId" env-required:"true"`
	Topics      []string `yaml:"topics" env-separator:"," env-required:"true"`
	StartOffset string   `yaml:"startOffset" env-default:"latest"` // 'latest' или 'earliest'

	MinBytes int           `yaml:"minBytes" env-default:"10000"`    // 10KB - Минимальный размер пакета для Fetch
	MaxBytes int           `yaml:"maxBytes" env-default:"10000000"` // 10MB - Максимальный размер пакета для Fetch
	MaxWait  time.Duration `yaml:"maxWait" env-default:"1s"`        // Макс. время ожидания MinBytes

	CommitInterval    time.Duration `yaml:"commitInterval" env-default:"1s"`    // Интервал авто-коммита (0 - отключает авто-коммит)
	HeartbeatInterval time.Duration `yaml:"heartbeatInterval" env-default:"3s"` // Частота отправки heartbeat брокеру
	SessionTimeout    time.Duration `yaml:"sessionTimeout" env-default:"30s"`   // Таймаут сессии консьюмера
	RebalanceTimeout  time.Duration `yaml:"rebalanceTimeout" env-default:"60s"` // Таймаут для ребалансировки

	DialTimeout  time.Duration `yaml:"dialTimeout" env-default:"3s"`   // Таймаут подключения к брокеру
	ReadTimeout  time.Duration `yaml:"readTimeout" env-default:"30s"`  // Таймаут чтения сообщений
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"10s"` // Таймаут записи (для коммитов и т.д.)
	MaxAttempts  int           `yaml:"maxAttempts" env-default:"3"`    // Макс. кол-во попыток для некоторых операций
}

type KafkaTopics struct {
	List              []string `yaml:"list" env-separator:"," env-required:"true"`
	NumPartitions     int      `yaml:"numPartitions" env-required:"true"`
	ReplicationFactor int      `yaml:"replicationFactor" env-required:"true"`
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
