package utils

import (
	"log"

	"github.com/go-playground/locales"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var validate *validator.Validate
var trans ut.Translator
var english locales.Translator

type ErrorData struct {
	Field      []Field  `json:"field"`
	InputValue []any    `json:"input_value"`
	Error      []string `json:"error"`
}

type Field struct {
	Name         string `json:"name"`
	ErrorMessage string `json:"message"`
}

// IsValid validate the given struct based on its rule
func init() {
	validate = validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ = uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)
}

// IsValid validate the given struct based on its rule
func IsValid(data any) (ErrorData, bool) {
	funcName := "shared.validator.IsValid"
	if err := validate.Struct(data); err != nil {
		log.Printf("%v: Validation failed: %v", funcName, err)
		errs := translateError(err, trans)
		return errs, false
	}
	return ErrorData{}, true
}

func translateError(err error, trans ut.Translator) (errs ErrorData) {
	if err == nil {
		return ErrorData{}
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		field := Field{}
		field.Name = e.Field()
		field.ErrorMessage = e.Translate(trans)
		errs.Field = append(errs.Field, field)
		errs.InputValue = append(errs.InputValue, e.Value())
	}
	return errs
}
