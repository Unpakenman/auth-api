package config

import (
	"auth-api/internal/app/constants"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"os"
	"time"
)

type Values struct {
	HttpServer      *HTTPServerConfig      `envconfig:"HTTP" required:"true"`
	App             *AppConfig             `envconfig:"APP" required:"true"`
	ClinicsDB       *DBConfig              `envconfig:"DB_CLINICS" required:"true"`
	LogLevel        string                 `envconfig:"LOG_LEVEL" required:"true"`
	HttpClient      *HTTPClientConfig      `envconfig:"HTTP_CLIENT" required:"true"`
	SomeHttpService *SomeHttpServiceConfig `envconfig:"SOME_HTTP_SERVICE" required:"true"`
	Redis           *Redis                 `envconfig:"REDIS" required:"true"`
	AMQPServer      *AMQPConfig            `envconfig:"RABBITMQ" required:"true"`
	GRPCServer      *GRPCServerConfig      `envconfig:"GRPC_SERVER" required:"true"`
	GRPCClient      *GRPCClientConfig      `envconfig:"GRPC_CLIENT" required:"true"`
}

type GRPCServerConfig struct {
	Port                           int32         `envconfig:"PORT" required:"true"`
	KeepaliveMaxConnectionIdle     time.Duration `envconfig:"KEEPALIVE_MAX_CONNECTION_IDLE" required:"true"`
	KeepaliveMaxConnectionAge      time.Duration `envconfig:"KEEPALIVE_MAX_CONNECTION_AGE" required:"true"`
	KeepaliveMaxConnectionAgeGrace time.Duration `envconfig:"KEEPALIVE_MAX_CONNECTION_AGE_GRACE" required:"true"`
	KeepaliveTime                  time.Duration `envconfig:"KEEPALIVE_TIME" required:"true"`
	KeepaliveTimeout               time.Duration `envconfig:"KEEPALIVE_TIMEOUT" required:"true"`
}

type GRPCClientConfig struct {
	KeepaliveTime               time.Duration `envconfig:"KEEPALIVE_TIME" required:"true"`
	KeepaliveTimeout            time.Duration `envconfig:"KEEPALIVE_TIMEOUT" required:"true"`
	KeepalivePermitWithoutCalls bool          `envconfig:"KEEPALIVE_PERMIT_WITHOUT_CALLS" required:"true"`
}

type HTTPServerConfig struct {
	ServerPort     int32  `envconfig:"HTTP_SERVER_PORT" required:"true"`
	ApiDefaultPath string `envconfig:"HTTP_API_DEFAULT_PATH" required:"true"`
}

type AppConfig struct {
	Name    string `envconfig:"NAME" required:"true"`
	Build   string `envconfig:"BUILD" required:"true"`
	SiteURL string `envconfig:"SITE_URL" required:"false"`
}

type DBConfig struct {
	DBName   string `envconfig:"DB_CLINICS_NAME" required:"true"`
	User     string `envconfig:"DB_CLINICS_USER" required:"true"`
	Password string `envconfig:"DB_CLINICS_PASSWORD" required:"true"`
	Hostname string `envconfig:"DB_CLINICS_HOSTNAME" required:"true"`
	SSLMode  string `envconfig:"DB_CLINICS_SSLMODE" required:"false"`
	Port     int32  `envconfig:"DB_CLINICS_PORT" required:"true"`

	MaxOpenConns                    int           `envconfig:"DB_CLINICS_MAX_OPEN_CONNS" required:"true"`
	MaxIdleConns                    int           `envconfig:"DB_CLINICS_MAX_IDLE_CONNS" required:"true"`
	MaxLifeTimeConns                time.Duration `envconfig:"DB_CLINICS_MAX_LIFETIME_CONNS" required:"true"`
	StatementTimeout                time.Duration `envconfig:"DB_CLINICS_STATEMENT_TIMEOUT" required:"false"`
	IdleInTransactionSessionTimeout time.Duration `envconfig:"DB_CLINICS_IDLE_IN_TRANSACTION_SESSION_TIMEOUT" required:"false"`
	LockTimeout                     time.Duration `envconfig:"DB_CLINICS_LOCK_TIMEOUT" required:"false"`
}

type Redis struct {
	URL    string `envconfig:"URL" required:"true"`
	Prefix string `envconfig:"PREFIX" required:"true"`
}

type AMQPConfig struct {
	User     string `envconfig:"USER" required:"true"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Hostname string `envconfig:"HOSTNAME" required:"true"`
	Protocol string `envconfig:"PROTOCOL" required:"true"`
	VHost    string `envconfig:"VHOST" required:"false"`
	Port     int32  `envconfig:"PORT" required:"true"`

	ReconnectTimeout time.Duration `envconfig:"RECONNECT_TIMEOUT_MS" required:"true"`

	EventsQueue           string `envconfig:"EVENTS_QUEUE" required:"true"`
	EventsPrefetchCount   int    `envconfig:"EVENTS_PREFETCH_COUNT" required:"true"`
	CommandsQueue         string `envconfig:"COMMANDS_QUEUE" required:"true"`
	CommandsPrefetchCount int    `envconfig:"COMMANDS_PREFETCH_COUNT" required:"true"`
	SmsQueue              string `envconfig:"SMS_QUEUE" required:"true"`
	SmsSender             string `envconfig:"SMS_SENDER" required:"true"`

	AppointmentsCommandsExchange string `envconfig:"APPOINTMENTS_EXCHANGE" required:"true"`
}

type HTTPClientConfig struct {
	Timeout time.Duration `envconfig:"TIMEOUT" required:"true"`
}

type SomeHttpServiceConfig struct {
	URL string `envconfig:"URL" required:"true"`
}

type flagsValues struct {
	Mode            string
	JobName         string
	UseLocalEnvFile bool
}

var Config *Values

func New() (*Values, error) {
	err := LoadEnvFile()
	if err != nil {
		return nil, err
	}

	cfg := &Values{}
	err = envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadEnvFile() error {
	if needUseLocalEnvFile() {
		err := godotenv.Load(constants.DefaultEnvFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func needUseLocalEnvFile() bool {
	for _, arg := range os.Args {
		if arg == constants.UseLocalEnvFileArg {
			return true
		}
	}
	return false
}
