// Package v1 representa tipos usados pela aplicação na versão 1
package v1

import (
	"errors"
)

// ErrorResponse é a forma usada para respostas da API sobre erros na API
type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// RequestError é usada para passar o erro durante a requisição através da
// aplicação com um contexto específico
type RequestError struct {
	Err    error
	Status int
}

// NewRequestError recebe um erro e o engloba com um status HTTP
// Deve ser usada quando handlers encontram erros esperados
func NewRequestError(err error, status int) error {
	return &RequestError{err, status}
}

// Error implementa a interface error usando a mensagem padrão do erro encontrado.
// É o que será mostrado nos logs da aplicação
func (re *RequestError) Error() string {
	return re.Err.Error()
}

// IsRequestError
func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

// GetRequestError retorna uma cópia do ponteiro do RequestError
func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}
