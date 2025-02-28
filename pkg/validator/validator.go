package validator

import (
	"errors"
	"reflect"
	"strings"
	"sync"

	enlocales "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"

	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

type IValidator interface {
	ValidateStruct(data interface{}) error
	ValidateVariable(data interface{}, tag string) error
}

type validatorStruct struct {
	validator  *validator.Validate
	translator ut.Translator
}

var (
	validatorInstance IValidator
	once              sync.Once
)

func NewValidator() IValidator {
	once.Do(func() {
		en := enlocales.New()
		translator := ut.New(en, en)

		trans, found := translator.GetTranslator("en")
		if !found {
			log.Error(map[string]interface{}{
				"error": "translator not found",
			}, "[VALIDATOR][NewValidator] Translator not found")
		}

		val := validator.New()

		val.RegisterTagNameFunc(func(fld reflect.StructField) string {
			for _, tagName := range []string{"json", "form", "query", "param"} {
				if tag := fld.Tag.Get(tagName); tag != "" {
					name := strings.SplitN(tag, ",", 2)[0]
					if name == "-" {
						continue
					}
					return name
				}
			}
			// Fall back to field name if no tags are found
			return fld.Name
		})

		err := entranslations.RegisterDefaultTranslations(val, trans)
		if err != nil {
			log.Error(map[string]interface{}{
				"error": err.Error(),
			}, "[VALIDATOR][NewValidator] Failed to register default translations")
		}

		validatorInstance = &validatorStruct{
			validator:  val,
			translator: trans,
		}
	})

	return validatorInstance
}

func (v *validatorStruct) ValidateStruct(data interface{}) error {
	if err := v.validator.Struct(data); err != nil {
		return v.handleValidationErrors(err, data)
	}
	return nil
}

func (v *validatorStruct) ValidateVariable(data interface{}, tag string) error {
	if err := v.validator.Var(data, tag); err != nil {
		return v.handleValidationErrors(err, nil)
	}
	return nil
}

func (v *validatorStruct) handleValidationErrors(err error, data interface{}) error {
	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		length := len(valErrs)
		resp := make(ValidationErrors, length)
		for i, err := range valErrs {
			fieldPath := err.Namespace()

			parts := strings.Split(fieldPath, ".")
			if len(parts) > 1 {
				parts = parts[1:]
			}

			jsonTag := strings.ToLower(strings.Join(parts, "."))

			resp[i] = map[string]validationError{
				jsonTag: {
					Tag:         err.Tag(),
					Param:       err.Param(),
					Translation: err.Translate(v.translator),
				},
			}
		}
		return resp
	}

	log.Error(map[string]interface{}{
		"error": err.Error(),
	}, "[VALIDATOR][handleValidationErrors] Unexpected error")
	return err
}
