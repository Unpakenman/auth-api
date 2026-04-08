package log

import (
	"github.com/sirupsen/logrus"
	"github.com/wI2L/jettison"
)

type ConsoleOptions struct {
	Level string
}

func NewConsole(opts ConsoleOptions) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.SetFormatter(&JSONFormatter{})
	level, err := logrus.ParseLevel(opts.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)

	return logger, nil
}

func formatConsoleExtrasWithObscure(sensitiveFields []string, value interface{}) []byte {
	rawJSONValue, err := jettison.MarshalOpts(value, jettison.NoHTMLEscaping(), jettison.DenyList(sensitiveFields))
	if err != nil {
		return []byte("failed serialize value to json string")
	}
	return rawJSONValue
}

func NewFields(key string, value ...interface{}) logrus.Fields {
	return logrus.Fields{
		key: value,
	}
}
