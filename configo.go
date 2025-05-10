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
	ProxyUrl  string `yaml:"proxyUrl"`
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

type Rest struct {
	Host               string                 `yaml:"host" env-default:"0.0.0.0"`                // Хост сервера
	Port               string                 `yaml:"port" env-default:"8080"`                   // Порт сервера
	ReadTimeout        time.Duration          `yaml:"readTimeout" env-default:"10s"`             // Таймаут чтения всего запроса
	WriteTimeout       time.Duration          `yaml:"writeTimeout" env-default:"10s"`            // Таймаут записи всего ответа
	IdleTimeout        time.Duration          `yaml:"idleTimeout" env-default:"60s"`             // Таймаут простоя keep-alive соединения
	HandlerTimeout     time.Duration          `yaml:"handlerTimeout" env-default:"15s"`          // Таймаут на обработку одного запроса (для middleware.Timeout)
	ShutdownTimeout    time.Duration          `yaml:"shutdownTimeout" env-default:"15s"`         // Таймаут на корректное завершение работы
	BaseURL            string                 `yaml:"baseURL"`                                   // Полный базовый URL сервера (для генерации ссылок)
	BasePath           string                 `yaml:"basePath" env-default:"/"`                  // Базовый путь для всех маршрутов API (например, "/api/v1")
	MaxRequestBodySize int64                  `yaml:"maxRequestBodySize" env-default:"10485760"` // Максимальный размер тела запроса в байтах (10MB)
	Compression        RestCompression        `yaml:"compression"`
	CORS               RestCORS               `yaml:"cors"`
	TLS                RestTLS                `yaml:"tls"`
	Profiling          RestProfiling          `yaml:"profiling"`
	RateLimit          RestRateLimitConfig    `yaml:"rateLimit"`
	SecurityHeaders    RestSecurityHeaders    `yaml:"securityHeaders"`
	StaticFiles        []RestFilesConfigEntry `yaml:"staticFiles"`                      // Массив для конфигурации раздачи нескольких наборов статики
	TrustedProxies     []string               `yaml:"trustedProxies" env-separator:","` // Список CIDR доверенных прокси
}

type RestCompression struct {
	Enabled bool `yaml:"enabled" env-default:"false"`
	Level   int  `yaml:"level" env-default:"-1"` // -1 соответствует flate.DefaultCompression
}

type RestCORS struct {
	Enabled            bool          `yaml:"enabled" env-default:"false"`
	AllowedOrigins     []string      `yaml:"allowedOrigins" env-separator:"," env-default:""` // По умолчанию пусто, можно установить ["*"] если нужно разрешить всем
	AllowedMethods     []string      `yaml:"allowedMethods" env-separator:"," env-default:"GET,POST,PUT,DELETE,OPTIONS"`
	AllowedHeaders     []string      `yaml:"allowedHeaders" env-separator:"," env-default:"Origin,Content-Type,Accept,Authorization"`
	ExposedHeaders     []string      `yaml:"exposedHeaders" env-separator:"," env-default:""`
	AllowCredentials   bool          `yaml:"allowCredentials" env-default:"false"`
	MaxAge             time.Duration `yaml:"maxAge" env-default:"300s"` // 5 минут
	OptionsPassthrough bool          `yaml:"optionsPassthrough" env-default:"false"`
	Debug              bool          `yaml:"debug" env-default:"false"`
}

type RestTLS struct {
	Enabled               bool     `yaml:"enabled" env-default:"false"`
	CertFile              string   `yaml:"certFile"`
	KeyFile               string   `yaml:"keyFile"`
	AutoCert              bool     `yaml:"autoCert" env-default:"false"`
	AutoCertCacheDir      string   `yaml:"autoCertCacheDir" env-default:".autocert"`
	AutoCertHostWhitelist []string `yaml:"autoCertHostWhitelist" env-separator:"," env-default:""`
}

type RestProfiling struct {
	Enabled bool   `yaml:"enabled" env-default:"false"`
	Prefix  string `yaml:"prefix" env-default:"/debug/pprof"`
}

type RestRateLimitConfig struct {
	Enabled         bool          `yaml:"enabled" env-default:"false"`
	RPS             float64       `yaml:"rps" env-default:"100"`
	Burst           int           `yaml:"burst" env-default:"20"`
	CleanupInterval time.Duration `yaml:"cleanupInterval" env-default:"1m"`
}

type RestSecurityHeaders struct {
	Enabled               bool   `yaml:"enabled" env-default:"true"`
	HSTSMaxAgeSeconds     int    `yaml:"hstsMaxAgeSeconds" env-default:"31536000"` // 1 год
	HSTSIncludeSubdomains bool   `yaml:"hstsIncludeSubdomains" env-default:"true"`
	HSTSPreload           bool   `yaml:"hstsPreload" env-default:"false"`
	ContentTypeNosniff    bool   `yaml:"contentTypeNosniff" env-default:"true"`
	FrameOptions          string `yaml:"frameOptions" env-default:"SAMEORIGIN"` // "DENY" или "SAMEORIGIN"
	XSSProtection         string `yaml:"xssProtection" env-default:"0"`         // "0" (CSP предпочтительнее), "1", "1; mode=block"
	ContentSecurityPolicy string `yaml:"contentSecurityPolicy" env-default:"default-src 'self'"`
	ReferrerPolicy        string `yaml:"referrerPolicy" env-default:"strict-origin-when-cross-origin"`
	PermissionsPolicy     string `yaml:"permissionsPolicy" env-default:""` // Пример: "geolocation=(), microphone=()"
}

type RestFilesConfigEntry struct {
	Enabled   bool          `yaml:"enabled" env-default:"false"`
	URLPrefix string        `yaml:"urlPrefix"` // Должен заканчиваться на "/"
	FSRoot    string        `yaml:"fsRoot"`
	CacheTTL  time.Duration `yaml:"cacheTtl" env-default:"1h"`
	IndexFile string        `yaml:"indexFile" env-default:"index.html"`
	SPA       bool          `yaml:"spa" env-default:"false"`
}

type GrpcServer struct {
	// Адрес для прослушивания
	Host string `yaml:"host" env:"GRPC_HOST" env-default:"0.0.0.0"`        // Хост, на котором сервер будет слушать
	Port int    `yaml:"port" env:"GRPC_PORT,required" env-required:"true"` // Порт для прослушивания

	// Конфигурация TLS
	EnableTLS    bool   `yaml:"enableTLS" env:"GRPC_ENABLE_TLS" env-default:"false"` // Включить ли TLS
	CertFile     string `yaml:"certFile" env:"GRPC_CERT_FILE"`                       // Путь к файлу сертификата сервера
	KeyFile      string `yaml:"keyFile" env:"GRPC_KEY_FILE"`                         // Путь к файлу приватного ключа сервера
	ClientCAFile string `yaml:"clientCAFile" env:"GRPC_CLIENT_CA_FILE"`              // Путь к файлу CA сертификата клиента для mTLS (опционально)

	// KeepAlive параметры сервера (grpc.KeepaliveServerParameters - KASP)
	// Управляют поведением сервера относительно keep-alive от клиентов и инициируемых сервером пингов.
	// Если значение Duration равно 0, часто используется значение по умолчанию gRPC (бесконечность или очень большое).
	KeepAliveMaxConnectionIdle     time.Duration `yaml:"keepAliveMaxConnectionIdle" env:"GRPC_KEEP_ALIVE_MAX_CONNECTION_IDLE" env-default:"0s"`          // Пример: "30m". Если клиент неактивен это время, сервер посылает GOAWAY. 0s = бесконечность (gRPC default).
	KeepAliveMaxConnectionAge      time.Duration `yaml:"keepAliveMaxConnectionAge" env:"GRPC_KEEP_ALIVE_MAX_CONNECTION_AGE" env-default:"0s"`            // Пример: "2h". Максимальное время жизни соединения. 0s = бесконечность (gRPC default).
	KeepAliveMaxConnectionAgeGrace time.Duration `yaml:"keepAliveMaxConnectionAgeGrace" env:"GRPC_KEEP_ALIVE_MAX_CONNECTION_AGE_GRACE" env-default:"0s"` // Пример: "30s". Дополнительное время после MaxConnectionAge. 0s = бесконечность (gRPC default).
	KeepAliveServerTime            time.Duration `yaml:"keepAliveServerTime" env:"GRPC_KEEP_ALIVE_SERVER_TIME" env-default:"2h"`                         // Если клиент неактивен это время, сервер пингует клиента. (gRPC default: 2h)
	KeepAliveServerTimeout         time.Duration `yaml:"keepAliveServerTimeout" env:"GRPC_KEEP_ALIVE_SERVER_TIMEOUT" env-default:"20s"`                  // Время ожидания ответа на пинг перед закрытием соединения. (gRPC default: 20s)

	// Политика принудительного KeepAlive (grpc.keepalive.EnforcementPolicy - KAEP)
	// Защищает сервер от злоупотреблений настройками keep-alive со стороны клиента.
	KeepAliveEnforcementPolicyMinTime             time.Duration `yaml:"keepAliveEnforcementPolicyMinTime" env:"GRPC_KEEP_ALIVE_ENFORCEMENT_MIN_TIME" env-default:"5m"`                             // Минимальный интервал между пингами от клиента. (gRPC default: 5m)
	KeepAliveEnforcementPolicyPermitWithoutStream bool          `yaml:"keepAliveEnforcementPolicyPermitWithoutStream" env:"GRPC_KEEP_ALIVE_ENFORCEMENT_PERMIT_WITHOUT_STREAM" env-default:"false"` // Разрешать ли пинги без активных потоков. (gRPC default: false)

	// Лимиты размеров сообщений
	MaxReceiveMessageSize int `yaml:"maxReceiveMessageSize" env:"GRPC_MAX_RECEIVE_MESSAGE_SIZE" env-default:"4194304"` // 4MB (gRPC default: 4MB)
	MaxSendMessageSize    int `yaml:"maxSendMessageSize" env:"GRPC_MAX_SEND_MESSAGE_SIZE" env-default:"0"`             // 0 для использования gRPC default (math.MaxInt32)

	// Лимиты потоков и параллелизма
	MaxConcurrentStreams  uint32 `yaml:"maxConcurrentStreams" env:"GRPC_MAX_CONCURRENT_STREAMS" env-default:"0"`    // Максимальное количество одновременных потоков на одном HTTP/2 соединении. 0 для gRPC default (math.MaxUint32).
	InitialWindowSize     int32  `yaml:"initialWindowSize" env:"GRPC_INITIAL_WINDOW_SIZE" env-default:"0"`          // Начальный размер окна для потока. 0 для gRPC default (64KB).
	InitialConnWindowSize int32  `yaml:"initialConnWindowSize" env:"GRPC_INITIAL_CONN_WINDOW_SIZE" env-default:"0"` // Начальный размер окна для соединения. 0 для gRPC default.

	// Размеры буферов (для настройки производительности)
	ReadBufferSize  int `yaml:"readBufferSize" env:"GRPC_READ_BUFFER_SIZE" env-default:"32768"`   // 32KB (gRPC default)
	WriteBufferSize int `yaml:"writeBufferSize" env:"GRPC_WRITE_BUFFER_SIZE" env-default:"32768"` // 32KB (gRPC default)

	// Регистрация стандартных сервисов
	EnableHealthCheckService bool `yaml:"enableHealthCheckService" env:"GRPC_ENABLE_HEALTH_CHECK_SERVICE" env-default:"true"` // Автоматически регистрировать стандартный Health Check сервис.
	EnableReflectionService  bool `yaml:"enableReflectionService" env:"GRPC_ENABLE_REFLECTION_SERVICE" env-default:"true"`    // Автоматически регистрировать Reflection сервис (для grpcurl, etc.).

	// "Вежливое" завершение работы (Graceful Shutdown)
	GracefulShutdownTimeout time.Duration `yaml:"gracefulShutdownTimeout" env:"GRPC_GRACEFUL_SHUTDOWN_TIMEOUT" env-default:"30s"` // Таймаут для ожидания завершения активных RPC перед принудительной остановкой.

	// Таймаут на установку соединения (до TLS handshake и создания потоков)
	ConnectionTimeout time.Duration `yaml:"connectionTimeout" env:"GRPC_CONNECTION_TIMEOUT" env-default:"120s"` // (gRPC default: 120s)

}

type GrpcClient struct {
	Host      string `yaml:"host" env:"HOST,required"`
	Port      int    `yaml:"port" env:"PORT,required"`
	UserAgent string `yaml:"userAgent" env:"USER_AGENT"`

	// TLS
	EnableTLS          bool   `yaml:"enableTLS" env:"ENABLE_TLS" env-default:"false"`
	CACertFile         string `yaml:"caCertFile" env:"CA_CERT_FILE"`                 // Путь к CA сертификату сервера
	ClientCertFile     string `yaml:"clientCertFile" env:"CLIENT_CERT_FILE"`         // Путь к клиентскому сертификату (для mTLS)
	ClientKeyFile      string `yaml:"clientKeyFile" env:"CLIENT_KEY_FILE"`           // Путь к приватному ключу клиента (для mTLS)
	ServerNameOverride string `yaml:"serverNameOverride" env:"SERVER_NAME_OVERRIDE"` // Переопределение имени сервера в TLS сертификате

	// Параметры повторных попыток подключения (для метода Connect)
	ConnectMaxAttempts       int           `yaml:"connectMaxAttempts" env:"CONNECT_MAX_ATTEMPTS" env-default:"5"`
	ConnectInitialBackoff    time.Duration `yaml:"connectInitialBackoff" env:"CONNECT_INITIAL_BACKOFF" env-default:"250ms"`
	ConnectMaxBackoff        time.Duration `yaml:"connectMaxBackoff" env:"CONNECT_MAX_BACKOFF" env-default:"5s"`
	ConnectBackoffMultiplier float64       `yaml:"connectBackoffMultiplier" env:"CONNECT_BACKOFF_MULTIPLIER" env-default:"2.0"`

	// Тайм-аут на установку соединения (для каждой отдельной попытки grpc.DialContext)
	DialTimeout time.Duration `yaml:"dialTimeout" env:"DIAL_TIMEOUT" env-default:"5s"`

	// KeepAlive параметры
	KeepAliveTime       time.Duration `yaml:"keepAliveTime" env:"KEEP_ALIVE_TIME" env-default:"30s"`
	KeepAliveTimeout    time.Duration `yaml:"keepAliveTimeout" env:"KEEP_ALIVE_TIMEOUT" env-default:"20s"`
	PermitWithoutStream bool          `yaml:"permitWithoutStream" env:"PERMIT_WITHOUT_STREAM" env-default:"true"`

	// Размеры сообщений
	MaxRecvMsgSize int `yaml:"maxRecvMsgSize" env:"MAX_RECV_MSG_SIZE" env-default:"4194304"` // 4MB
	MaxSendMsgSize int `yaml:"maxSendMsgSize" env:"MAX_SEND_MSG_SIZE" env-default:"4194304"` // 4MB
}

type Grpc struct {
	Server  GrpcServer   `yaml:"server"`
	Clients []GrpcClient `yaml:"clients"`
}

type Ws struct {
	// Включен ли WebSocket сервер
	Enabled bool `yaml:"enabled" env-default:"true"`
	// Хост, на котором будет слушать WebSocket сервер (например, "0.0.0.0" для всех интерфейсов)
	Host string `yaml:"host" env-default:"0.0.0.0"`
	// Порт, на котором будет слушать WebSocket сервер
	Port int `yaml:"port" env-required:"true"`
	// Путь эндпоинта для WebSocket соединений (например, "/ws")
	Path string `yaml:"path" env-default:"/ws"`

	// TLS конфигурация
	EnableTLS bool   `yaml:"enableTLS" env-default:"false"`
	CertFile  string `yaml:"certFile" env:"WS_CERT_FILE"` // Путь к файлу SSL сертификата
	KeyFile   string `yaml:"keyFile" env:"WS_KEY_FILE"`   // Путь к файлу приватного ключа SSL

	// Параметры соединения
	// Таймаут для WebSocket handshake
	HandshakeTimeout time.Duration `yaml:"handshakeTimeout" env-default:"5s"`
	// Размер буфера чтения для каждого соединения (в байтах)
	ReadBufferSize int `yaml:"readBufferSize" env-default:"4096"`
	// Размер буфера записи для каждого соединения (в байтах)
	WriteBufferSize int `yaml:"writeBufferSize" env-default:"4096"`
	// Максимальный размер одного входящего сообщения (в байтах)
	MaxMessageReadSize int64 `yaml:"maxMessageReadSize" env-default:"65536"` // 64KB

	// Сжатие (permessage-deflate)
	// Включить сжатие сообщений
	EnableCompression bool `yaml:"enableCompression" env-default:"false"`
	// Уровень сжатия (если включено). -1 для значения по умолчанию (обычно это 6), 0 - без сжатия, 1-9 - уровни сжатия.
	CompressionLevel int `yaml:"compressionLevel" env-default:"-1"`

	// Безопасность и лимиты
	// Список разрешенных origins для CORS. Пустой список или ["*"] для разрешения всех.
	// Пример: "http://localhost:3000,https://example.com"
	AllowedOrigins []string `yaml:"allowedOrigins" env-separator:"," env-default:""`
	// Поддерживаемые WebSocket субпротоколы. Клиент может запросить один из них.
	// Пример: "chat,json-rpc"
	Subprotocols []string `yaml:"subprotocols" env-separator:"," env-default:""`
	// Максимальное общее количество активных WebSocket соединений (0 - без ограничений)
	MaxConnections int `yaml:"maxConnections" env-default:"0"`
	// Максимальное количество соединений с одного IP-адреса
	MaxConnectionsPerIP int `yaml:"maxConnectionsPerIP" env-required:"true"`

	// Корректное завершение работы
	// Таймаут для ожидания завершения активных соединений перед принудительной остановкой сервера
	ShutdownTimeout time.Duration `yaml:"shutdownTimeout" env-default:"10s"`

	// Конфигурация сессии WebSocket
	Session WsSession `yaml:"session"`
}

// WsSession конфигурация для индивидуальной WebSocket сессии
type WsSession struct {
	// Ping/Pong для поддержания активности соединения (инициируется сервером)
	// Включить отправку ping-сообщений сервером
	EnablePing bool `yaml:"enablePing" env-default:"true"`
	// Интервал, с которым сервер отправляет ping-сообщения клиенту, если от клиента нет активности.
	PingInterval time.Duration `yaml:"pingInterval" env-default:"30s"`
	// Таймаут ожидания pong-сообщения от клиента после отправки ping. Если pong не получен, соединение закрывается.
	PongTimeout time.Duration `yaml:"pongTimeout" env-default:"10s"`

	// Максимальное время простоя соединения (без обмена данными, включая контрольные фреймы, если они неактивны или нечасты).
	// Сервер может закрыть соединение, если оно простаивает дольше этого времени. 0 - отключить проверку.
	MaxIdleTime time.Duration `yaml:"maxIdleTime" env-default:"60s"`

	// Буферизация исходящих сообщений (от сервера к клиенту)
	// Размер буфера (количество сообщений) для исходящих сообщений для одного клиента.
	OutboundMessageBufferSize int `yaml:"outboundMessageBufferSize" env-default:"256"`
	// Таймаут на запись одного сообщения клиенту. Предотвращает блокировку, если клиент медленно читает.
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"5s"`

	// Максимальное время жизни сессии
	// Максимальная продолжительность существования WebSocket сессии, независимо от активности.
	// 0 - без ограничения времени жизни.
	MaxLifetime time.Duration `yaml:"maxLifetime" env-default:"0s"`
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
