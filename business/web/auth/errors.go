package auth

import (
	"errors"
	"fmt"
)

// AuthError é usado para passar erros durante a requisição que estejam envolvidos
// com o context de autenticação
type AuthError struct {
	msg string
}

// NewAuthError cria um AuthError com a mensagem definida
func NewAuthError(format string, args ...any) error {
	return &AuthError{
		msg: fmt.Sprintf(format, args...),
	}
}

// Error implementa a interface error. Usa a mensagem padrão do erro envolvido.
// é o que será mostrado nos logs
func (ae *AuthError) Error() string {
	return ae.msg
}

// IsAuthError checa se um erro é do tipo AuthError
func IsAuthError(err error) bool {
	var ae *AuthError
	return errors.As(err, &ae)
}
