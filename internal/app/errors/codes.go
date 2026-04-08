package errors

import (
	"net/http"
	"strconv"

	traceCodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
)

type Codes struct {
	HTTP uint32
	GRPC codes.Code
}

type StatusName string

const (
	StatusInternalServerError StatusName = "INTERNAL"
	StatusOK                  StatusName = "OK"
	StatusBadRequest          StatusName = "BAD_REQUEST"
	StatusNotFound            StatusName = "NOT_FOUND"
	StatusInvalidArgument     StatusName = "INVALID_ARGUMENT"
)

var statusCodesMap = map[StatusName]Codes{
	StatusInternalServerError: {
		HTTP: http.StatusInternalServerError,
		GRPC: codes.Unknown,
	},
	StatusOK: {
		HTTP: http.StatusOK,
		GRPC: codes.OK,
	},
	StatusBadRequest: {
		HTTP: http.StatusBadRequest,
		GRPC: codes.FailedPrecondition,
	},
	StatusNotFound: {
		HTTP: http.StatusNotFound,
		GRPC: codes.NotFound,
	},
	StatusInvalidArgument: {
		HTTP: http.StatusBadRequest,
		GRPC: codes.InvalidArgument,
	},
}

func GetCodesByStatusName(name StatusName) Codes {
	result, ok := statusCodesMap[name]
	if !ok {
		return Codes{}
	}
	return result
}

func HTTPStatusStringFromCode(code codes.Code) string {
	return strconv.Itoa(HTTPStatusFromCode(code))
}

func HTTPStatusFromCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return http.StatusInternalServerError
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.FailedPrecondition:
		return http.StatusBadRequest
	case codes.Aborted:
		return http.StatusConflict
	case codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func TraceCodeFromGRPCCode(code codes.Code) traceCodes.Code {
	if code == codes.OK {
		return traceCodes.Ok
	}
	return traceCodes.Error
}
func TraceCodeFromHTTPode(code string) traceCodes.Code {
	if code == "200" {
		return traceCodes.Ok
	}
	return traceCodes.Error
}
