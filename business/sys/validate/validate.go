// Package validate contém o suporte para validar models
package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator ut.Translator

func init() {

	// Instantiate a validator.
	validate = validator.New()

	// Create a translator for english so the error messages are
	// more human-readable than technical.
	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// Register the english error messages for use.
	en_translations.RegisterDefaultTranslations(validate, translator)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Check valida o modelo passado de acordo com as tags declaradas nos cmapos do struct
// exemplo: adicionamos a tag min=3, vai verificar se o valor que foi passado
// obedece à regra de ter pelo menos 3 caracteres
func Check(val any) error {
	if err := validate.Struct(val); err != nil {

		// Usa type assertion para buscar o valor real do erro
		verrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// quando retornarmos, o middleware de erros checa se o valor é do tipo FieldErrors
		// e aí trabalha em cima dessa informação
		var fields FieldErrors
		for _, verror := range verrors {
			// constrói os erros que definimos com base no erro retornado pela validação
			field := FieldError{
				Field: verror.Field(),
				Err:   verror.Translate(translator),
			}
			fields = append(fields, field)
		}

		return fields
	}

	return nil
}
