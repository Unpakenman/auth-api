package constants

import "time"

const (
	UseLocalEnvFileArg = "--use-local-env"
	DefaultEnvFile     = "../.env"

	GRPCTraceIDField    = "traceparent"
	HTTPTraceHeaderName = "traceparent"

	RabbitMQLabel = "RabbitMQ"
)

const (
	UseLocalEnvFileDescription = "use env file for configuration"
	ModeNameArg                = "mode"
	ServiceMode                = "service"
	JobMode                    = "job"
	ModeArgDescription         = "Run mode:\n\tservice - run service\n\tjob - run job."

	JobNameArgName    = "job-name"
	JobArgDescription = "Job name:\n\t<name> - run <name> job"
	JobNameEmptyError = "job-name for mode=job required"
)

const (
	CachedTTLappointmentList = 300 * time.Second
)
