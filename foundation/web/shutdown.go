package web

import (
	"errors"
)

// shutdownError tipo usado para ajudar no desligamento  seguro da aplicação
type shutdownError struct {
	Message string
}

// NewShutdownError retorna um erro que causa o framework sinalizar um desligamento
func NewShutdownError(message string) error {
	return &shutdownError{message}
}

// Error implementação da interface error
func (se *shutdownError) Error() string {
	return se.Message
}

// IsShutdown checa se o erro em questão é um erro de shutdown ou não
func IsShutdown(err error) bool {
	var se *shutdownError      // semantica de ponteiro para structs
	return errors.As(err, &se) // verifica se é um erro do tipo passado e retorna uma cópia se for
}
