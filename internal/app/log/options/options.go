package options

type LoggerOptions struct {
	Extras    any
	Exception *ExceptionData
	TraceID   string
	Protocol  Protocol
	RequestID string
	Handler   string
	GamblerID int64
	Phone     string
}

type ExceptionData struct {
	Message    string
	Type       string
	Stacktrace string
}

type Protocol string

const (
	HTTPProtocol Protocol = "http"
	GRPCProtocol Protocol = "grpc"
	AMQPProtocol Protocol = "amqp"
	NATSProtocol Protocol = "nats"
	CLIProtocol  Protocol = "cli"
)

type LoggerOption func(*LoggerOptions)

func WithException(data *ExceptionData) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.Exception = data
	}
}

func WithExtras(data any) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.Extras = data
	}
}

func WithTraceID(data string) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.TraceID = data
	}
}

func WithProtocol(data Protocol) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.Protocol = data
	}
}

func WithRequestID(data string) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.RequestID = data
	}
}

func WithHandler(data string) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.Handler = data
	}
}

func WithGamblerID(data int64) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.GamblerID = data
	}
}

func WithPhone(data string) LoggerOption {
	return func(lo *LoggerOptions) {
		lo.Phone = data
	}
}
