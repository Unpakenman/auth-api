package mapper

import (
	appErrors "auth-api/internal/app/errors"
	localerrors "auth-api/internal/app/errors"
	"errors"
	pbcommon "github.com/Unpakenman/proto/auth-api/gen/go/auth/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
)

func (m *mapper) toProtoError(code codes.Code, errorMessage string, details proto.Message) error {
	st := status.New(code, errorMessage)
	ds, err := st.WithDetails(protoadapt.MessageV1Of(details))
	if err != nil {
		return st.Err()
	}

	return ds.Err()
}

func (m *mapper) ResultErrorToProtoError(resultError localerrors.Error) error {
	errorMessage := resultError.Error()

	if errors.Is(resultError, appErrors.ErrInvalidDateTime) {
		details := &pbcommon.UserAlreadyExist{}
		return m.toProtoError(codes.AlreadyExists, errorMessage, details)
	}
	return resultError.ErrorProto()
}
