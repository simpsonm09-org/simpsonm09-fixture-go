package api

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// UseValidator installs the fleet request validator on Gin. It adds the
// `notblank` rule, which rejects a string that is empty or only whitespace.
func UseValidator() {
	binding.Validator = ginValidator{validate: newValidator()}
}

func newValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	_ = validate.RegisterValidation("notblank", notBlank)
	return validate
}

func notBlank(field validator.FieldLevel) bool {
	value, ok := field.Field().Interface().(string)
	if !ok {
		return false
	}
	return strings.TrimSpace(value) != ""
}

// ginValidator adapts *validator.Validate to Gin's binding.StructValidator.
type ginValidator struct {
	validate *validator.Validate
}

// ValidateStruct validates obj and returns a binding.SliceValidationError for a
// slice of structs, matching Gin's own validator behavior.
func (v ginValidator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}
	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Pointer:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validate.Struct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		errors := make(binding.SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				errors = append(errors, err)
			}
		}
		if len(errors) == 0 {
			return nil
		}
		return errors
	default:
		return nil
	}
}

// Engine returns the underlying validator engine.
func (v ginValidator) Engine() any {
	return v.validate
}
