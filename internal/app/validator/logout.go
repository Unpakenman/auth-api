package validator

import (
	localerrors "auth-api/internal/app/errors"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
	"github.com/gobuffalo/validate"
)

func (v *validator) Logout(req *pb.LogoutRequest) *[]localerrors.FieldViolation {
	checks := []validate.Validator{
		&StringLenGreaterThenValidator{
			Name:  "token",
			Field: req.Token,
			Min:   1,
		},
		&StringLenGreaterThenValidator{
			Name:  "refresh_token",
			Field: req.RefreshToken,
			Min:   1,
		},
	}
	errs := validate.Validate(checks...)
	return FormatValidateErrors(errs)
}
