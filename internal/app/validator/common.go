package validator

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"

	localerrors "auth-api/internal/app/errors"
	"github.com/gobuffalo/validate"
)

func FormatValidateErrors(errs *validate.Errors) *[]localerrors.FieldViolation {
	if !errs.HasAny() {
		return nil
	}
	formatted := []localerrors.FieldViolation{}
	for field, err := range errs.Errors {
		for _, errVal := range err {
			formatted = append(formatted, localerrors.FieldViolation{
				Field:       field,
				Description: errVal,
			})
		}
	}
	return &formatted
}

type UniqueSliceValidator[T comparable] struct {
	Name  string
	Field []T
}

func (v *UniqueSliceValidator[T]) IsValid(errors *validate.Errors) {
	keys := make(map[T]interface{})
	for _, entry := range v.Field {
		_, ok := keys[entry]
		if !ok {
			keys[entry] = nil
			continue
		}
		errors.Add(
			strings.ToLower(v.Name),
			"поле должно содержать только уникальные значения",
		)
		return
	}
}

type MinLenSliceValidator[T any] struct {
	Name  string
	Field []T
	Min   int
}

func (v *MinLenSliceValidator[T]) IsValid(errors *validate.Errors) {
	if len(v.Field) < v.Min {
		errors.Add(
			strings.ToLower(v.Name),
			fmt.Sprintf("элементов в поле должно быть больше чем %d", v.Min),
		)
	}
}

type numbers interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

type GreaterThenValueSliceValidator[T numbers] struct {
	Name  string
	Field []T
	Min   T
}

func (v *GreaterThenValueSliceValidator[T]) IsValid(errors *validate.Errors) {
	for _, entry := range v.Field {
		if entry <= v.Min {
			errors.Add(
				strings.ToLower(v.Name),
				fmt.Sprintf("элементы в поле должны быть больше чем %v", v.Min),
			)
			return
		}
	}
}

type IsGreaterThanValidator[T numbers] struct {
	Name  string
	Field T
	Min   T
}

func (v *IsGreaterThanValidator[T]) IsValid(errors *validate.Errors) {
	if v.Field > v.Min {
		return
	}

	errors.Add(
		strings.ToLower(v.Name),
		fmt.Sprintf("значение поля должно быть больше чем %v", v.Min),
	)
}

type ValueInListValidator[T comparable] struct {
	Name  string
	Field T
	List  []T
}

func (v *ValueInListValidator[T]) IsValid(errors *validate.Errors) {
	for _, validValue := range v.List {
		if v.Field == validValue {
			return
		}
	}
	errors.Add(
		strings.ToLower(v.Name),
		fmt.Sprintf("значение поля может быть только из списка %v", v.List),
	)
}

type SliceLenGreaterThenValidator[T any] struct {
	Name  string
	Field []T
	Min   int
}

func (v *SliceLenGreaterThenValidator[T]) IsValid(errors *validate.Errors) {
	if len(v.Field) > v.Min {
		return
	}
	errors.Add(
		strings.ToLower(v.Name),
		fmt.Sprintf("элементов в поле должно быть больше чем %v", v.Min),
	)
}

type StringLenGreaterThenValidator struct {
	Name    string
	Field   string
	Min     int
	Message string
}

func (v *StringLenGreaterThenValidator) IsValid(errors *validate.Errors) {
	strLength := utf8.RuneCountInString(v.Field)
	if v.Message == "" {
		v.Message = fmt.Sprintf("длина поля должна быть больше %d", v.Min)
	}
	if strLength <= v.Min {
		errors.Add(strings.ToLower(v.Name), v.Message)
	}
}

type OneOfNotNilSliceValidator struct {
	Name    string
	Field   []interface{}
	Message string
}

func (v *OneOfNotNilSliceValidator) IsValid(errors *validate.Errors) {
	counter := 0
	for _, val := range v.Field {
		if val != nil && !reflect.ValueOf(val).IsNil() {
			counter++
		}
	}
	if counter == 0 {
		errors.Add(
			strings.ToLower(v.Name),
			"хотя бы одно из полей должно быть заполнено",
		)
	}
}

type StringIsDateIsIso8601Validator struct {
	Name  string
	Field string
}

func (v *StringIsDateIsIso8601Validator) IsValid(errors *validate.Errors) {
	if v.Field == "" {
		return
	}
	if _, err := time.Parse(time.RFC3339, v.Field); err != nil {
		errors.Add(
			v.Name,
			"должен соответствовать формату ISO 8601",
		)
	}
}

type IsGreaterThanOrEqualValidator[T numbers] struct {
	Name  string
	Field T
	Min   T
}

func (v *IsGreaterThanOrEqualValidator[T]) IsValid(errors *validate.Errors) {
	if v.Field >= v.Min {
		return
	}

	errors.Add(
		v.Name,
		fmt.Sprintf("значение поля должно быть больше или равно %v", v.Min),
	)
}

type DateWithinMonthsInPastValidator struct {
	Name   string
	Field  string
	Months int
}

func (v *DateWithinMonthsInPastValidator) IsValid(errors *validate.Errors) {
	if v.Field == "" {
		return
	}
	fieldTime, err := time.Parse(time.RFC3339, v.Field)
	if err != nil {
		return
	}
	now := time.Now()
	minAllowedDate := now.AddDate(0, -v.Months, 0)
	if fieldTime.Before(minAllowedDate) {
		errors.Add(
			strings.ToLower(v.Name),
			fmt.Sprintf("дата не должна быть ранее чем %d месяцев назад от текущей даты", v.Months),
		)
	}
}

type DateIsGreaterThanValidator struct {
	Name  string
	Field string
	Min   string
}

func (v *DateIsGreaterThanValidator) IsValid(errors *validate.Errors) {
	if v.Field == "" {
		return
	}

	fieldTime, err := time.Parse(time.RFC3339, v.Field)
	if err != nil {
		return // Skip validation if field is not a valid date (handled by other validators)
	}

	minTime, err := time.Parse(time.RFC3339, v.Min)
	if err != nil {
		return // Skip validation if Min is not a valid date
	}

	if fieldTime.After(minTime) {
		return
	}

	errors.Add(
		strings.ToLower(v.Name),
		fmt.Sprintf("дата должна быть позже чем %s", v.Min),
	)
}

type IsLessThanOrEqualValidator[T numbers] struct {
	Name  string
	Field T
	Max   T
}

func (v *IsLessThanOrEqualValidator[T]) IsValid(errors *validate.Errors) {
	if v.Field <= v.Max {
		return
	}

	errors.Add(
		strings.ToLower(v.Name),
		fmt.Sprintf("значение поля должно быть меньше или равно %v", v.Max),
	)
}
