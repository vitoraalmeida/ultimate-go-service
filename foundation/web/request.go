package web

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

type validator interface {
	Validate() error
}

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	m := httptreemux.ContextParams(r.Context())
	return m[key]
}

// Decode lé o corpo de uma requisição HTTP buscando por um JSON.
// O JSON é decodificado para o objeto do valor (val) passado
// Se o valor for um struct, é checado de acordo com as tags de validação
// Se o valor possuir uma função de validação própria, é chamada
func Decode(r *http.Request, val any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	// checa se o valor passado é de um tipo que implementou um método Validate
	if v, ok := val.(validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}
