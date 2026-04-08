package log

import (
	"auth-api/internal/app/log/options"
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"

	"github.com/sirupsen/logrus"
)

type JSONFormatter struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	// The format to use is the same than for time.Format or time.Parse from the standard
	// library.
	// The standard Library already provides a set of predefined format.
	TimestampFormat string

	// DisableTimestamp allows disabling automatic timestamps in output
	DisableTimestamp bool

	// DisableHTMLEscape allows disabling html escaping in output
	DisableHTMLEscape bool

	// CallerPrettyfier can be set by the user to modify the content
	// of the function and file keys in the json data when ReportCaller is
	// activated. If any of the returned value is the empty string the
	// corresponding key will be removed from json fields.
	CallerPrettyfier func(*runtime.Frame) (function string, file string)

	// PrettyPrint will indent all json logs
	PrettyPrint bool
}

// LogDataFields uses for set order for fields in logs
type LogDataFields struct {
	Time    string          `json:"time"`
	Level   string          `json:"level"`
	Msg     LogDataMessage  `json:"msg,omitempty"`
	Context *LogDataContext `json:"context,omitempty"`
	Error   *LogDataError   `json:"error,omitempty"`
	Labels  *LogDataLabels  `json:"labels,omitempty"`
}

type LogDataError struct {
	Message    string `json:"message,omitempty"`
	Type       string `json:"type,omitempty"`
	Stacktrace string `json:"stacktrace,omitempty"`
}

type LogDataLabels struct {
	GamblerID int64  `json:"gambler_id,omitempty"`
	Phone     string `json:"phone,omitempty"`
}

type LogDataMessage struct {
	Message string `json:"message,omitempty"`
	Extras  any    `json:"extras,omitempty"`
	Func    string `json:"func,omitempty"`
	File    string `json:"file,omitempty"`
}

type LogDataContext struct {
	Protocol  string `json:"protocol,omitempty"`
	Handler   string `json:"handler,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	TraceID   string `json:"trace_id,omitempty"`
}

// Format renders a single log entry
func (f *JSONFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := LogDataFields{}
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data.Error = &LogDataError{
				Message: v.Error(),
			}

		default:
			if k == LogOptionsField {
				opts, ok := v.(options.LoggerOptions)
				if ok {
					showContext := opts.Handler != "" ||
						opts.RequestID != "" ||
						opts.Protocol != "" ||
						opts.TraceID != ""
					showLabels := opts.GamblerID != 0 || opts.Phone != ""
					showError := opts.Exception != nil
					showExtras := !(opts.Extras == nil ||
						reflect.ValueOf(opts.Extras).IsNil())
					if showContext {
						data.Context = &LogDataContext{}
						data.Context.Handler = opts.Handler
						data.Context.RequestID = opts.RequestID
						data.Context.Protocol = string(opts.Protocol)
						data.Context.TraceID = opts.TraceID
					}
					if showLabels {
						data.Labels = &LogDataLabels{}
						data.Labels.GamblerID = opts.GamblerID
						data.Labels.Phone = opts.Phone
					}
					if showError {
						data.Error = &LogDataError{
							Message:    opts.Exception.Message,
							Type:       opts.Exception.Type,
							Stacktrace: opts.Exception.Stacktrace,
						}
					}

					if showExtras {
						data.Msg.Extras = opts.Extras
					}
				}
			}
		}
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = DefaultTimestampFormat
	}

	if !f.DisableTimestamp {
		data.Time = entry.Time.Format(timestampFormat)
	}
	data.Msg.Message = entry.Message
	data.Level = entry.Level.String()
	if entry.HasCaller() {
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}
		if funcVal != "" {
			data.Msg.Func = funcVal
		}
		if fileVal != "" {
			data.Msg.File = fileVal
		}
	}

	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	encoder := json.NewEncoder(b)
	encoder.SetEscapeHTML(!f.DisableHTMLEscape)
	if f.PrettyPrint {
		encoder.SetIndent("", "  ")
	}
	if err := encoder.Encode(data); err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %w", err)
	}

	return b.Bytes(), nil
}
