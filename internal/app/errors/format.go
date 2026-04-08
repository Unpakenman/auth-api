package errors

import (
	"bytes"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

type FieldViolation struct {
	Field       string
	Description string
}

func FormatValidateErrors(errs []FieldViolation) error {
	if len(errs) == 0 {
		return nil
	}
	b := new(bytes.Buffer)
	for _, value := range errs {
		fmt.Fprintf(b, "%s ", value)
	}
	return fmt.Errorf("errors: %s", b.String())
}

func FormatValidateErrorsToBadRequest(
	errs []FieldViolation,
	statusCode Codes,
	msg string,
) *grpcstatus.Status {
	details := make([]*errdetails.BadRequest_FieldViolation, len(errs))
	for i, err := range errs {
		details[i] = &errdetails.BadRequest_FieldViolation{
			Field:       err.Field,
			Description: err.Description,
		}
	}
	data := &errdetails.BadRequest{
		FieldViolations: details,
	}
	status := grpcstatus.New(statusCode.GRPC, msg)
	status, _ = status.WithDetails(data)
	return status
}

func FormatValidateErrorsToPreconditionFailure(
	errs []FieldViolation,
	statusCode Codes,
	msg string,
) *grpcstatus.Status {
	details := make([]*errdetails.PreconditionFailure_Violation, len(errs))
	for i, err := range errs {
		details[i] = &errdetails.PreconditionFailure_Violation{
			Subject:     err.Field,
			Description: err.Description,
		}
	}
	data := &errdetails.PreconditionFailure{
		Violations: details,
	}
	status := grpcstatus.New(statusCode.GRPC, msg)
	status, _ = status.WithDetails(data)
	return status
}

func FormatGRPCError(msg string, code codes.Code, message protoadapt.MessageV1) *grpcstatus.Status {
	status := grpcstatus.New(code, msg)
	status, _ = status.WithDetails(message)
	return status
}

func FormatError(err error, statusCode Codes) *grpcstatus.Status {
	formatted := grpcstatus.New(statusCode.GRPC, err.Error())
	return formatted
}
