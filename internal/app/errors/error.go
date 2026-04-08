package errors

import (
	"errors"

	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type Error interface {
	StatusCode() Codes
	Error() string
	ErrorProto() error
	Unwrap() error
}

func NewBadRequestErr(err error) Error {
	statusCode := GetCodesByStatusName(StatusBadRequest)
	return &BadRequestErr{
		OrigError:  err,
		statusCode: statusCode,
		ProtoError: FormatError(err, statusCode),
	}
}

type BadRequestErr struct {
	OrigError  error
	statusCode Codes
	ProtoError *grpcstatus.Status
}

func (e *BadRequestErr) StatusCode() Codes {
	return e.statusCode
}

func (e *BadRequestErr) Error() string {
	return errors.Join(errors.New("bad request error: "), e.OrigError).Error()
}

func (e *BadRequestErr) Unwrap() error {
	return e.OrigError
}

func (e *BadRequestErr) ErrorProto() error {
	return e.ProtoError.Err()
}

func NewInternalErr(err error) Error {
	statusCode := GetCodesByStatusName(StatusInternalServerError)
	return &InternalErr{
		OrigError:  err,
		statusCode: statusCode,
		ProtoError: FormatError(err, statusCode),
	}
}

type InternalErr struct {
	OrigError  error
	statusCode Codes
	ProtoError *grpcstatus.Status
}

func (e *InternalErr) StatusCode() Codes {
	return e.statusCode
}

func (e *InternalErr) Error() string {
	return errors.Join(errors.New("internal error: "), e.OrigError).Error()
}

func (e *InternalErr) Unwrap() error {
	return e.OrigError
}

func (e *InternalErr) ErrorProto() error {
	return e.ProtoError.Err()
}

func NewNotFoundErr(err error) Error {
	statusCode := GetCodesByStatusName(StatusNotFound)
	return &InvalidArgumentErr{
		OrigError:  err,
		statusCode: statusCode,
		ProtoError: FormatError(err, statusCode),
	}
}

type NotFoundErr struct {
	OrigError  error
	statusCode Codes
	ProtoError *grpcstatus.Status
}

func (e *NotFoundErr) StatusCode() Codes {
	return e.statusCode
}

func (e *NotFoundErr) Error() string {
	return errors.Join(errors.New("not found error: "), e.OrigError).Error()
}

func (e *NotFoundErr) Unwrap() error {
	return e.OrigError
}

func (e *NotFoundErr) ErrorProto() error {
	return e.ProtoError.Err()
}

func NewInvalidArgumentErr(errs []FieldViolation) Error {
	statusCode := GetCodesByStatusName(StatusInvalidArgument)
	return &InvalidArgumentErr{
		OrigError:  FormatValidateErrors(errs),
		statusCode: statusCode,
		ProtoError: FormatValidateErrorsToBadRequest(
			errs,
			statusCode,
			"Ошибка валидации",
		),
	}
}

func NewProtoError(origError error, code codes.Code, msg string, message protoadapt.MessageV1) Error {
	return &InvalidArgumentErr{
		OrigError: origError,
		statusCode: Codes{
			HTTP: 0,
			GRPC: code,
		},
		ProtoError: FormatGRPCError(msg, code, message),
	}
}

type InvalidArgumentErr struct {
	OrigError  error
	statusCode Codes
	ProtoError *grpcstatus.Status
}

func (e *InvalidArgumentErr) StatusCode() Codes {
	return e.statusCode
}

func (e *InvalidArgumentErr) Error() string {
	return errors.Join(errors.New("ошибка валидации: "), e.OrigError).Error()
}

func (e *InvalidArgumentErr) Unwrap() error {
	return e.OrigError
}

func (e *InvalidArgumentErr) ErrorProto() error {
	return e.ProtoError.Err()
}
