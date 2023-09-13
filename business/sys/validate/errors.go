// erros de validação
package validate

import (
	"encoding/json"
	"errors"
)

// FieldError usado para indicar um erro com um campo específico de uma requisição
type FieldError struct {
	Field string `json:"field"`
	Err   string `json:"error"`
}

// FieldErrors representa uma coleção de fielErrors
type FieldErrors []FieldError

// NewFieldsError cria um FieldsError
func NewFieldsError(field string, err error) error {
	return FieldErrors{
		{
			Field: field,
			Err:   err.Error(),
		},
	}
}

// Error implementa a interface error
func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}

// Fields returna os campos que falharam na validação
func (fe FieldErrors) Fields() map[string]string {
	m := make(map[string]string)
	for _, fld := range fe {
		m[fld.Field] = fld.Err
	}
	return m
}

// IsFieldErrors checa se existe um erro do tipo FieldErrors
func IsFieldErrors(err error) bool {
	var fe FieldErrors
	return errors.As(err, &fe)
}

// GetFieldErrors returna uma cópia do ponteiro para um FieldErrors
func GetFieldErrors(err error) FieldErrors {
	var fe FieldErrors
	if !errors.As(err, &fe) {
		return nil
	}
	return fe
}
