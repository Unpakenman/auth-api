package validator

import (
	localerrors "auth-api/internal/app/errors"
	pb "github.com/Unpakenman/proto/auth-api/gen/go/auth/rpc"
	"github.com/gobuffalo/validate"
)

func (v *validator) Register(req *pb.RegisterRequest) *[]localerrors.FieldViolation {
	checks := []validate.Validator{
		&StringLenGreaterThenValidator{
			Name:  "phone",
			Field: req.PhoneNumber,
			Min:   10,
		},
		&StringLenGreaterThenValidator{
			Name:  "password",
			Field: req.Password,
			Min:   8,
		},
		&StringLenGreaterThenValidator{
			Name:  "email",
			Field: req.Email,
			Min:   8,
		},
	}
	errors := validate.Validate(checks...)
	return FormatValidateErrors(errors)
}
