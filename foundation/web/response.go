package web

import (
	"context"
	"encoding/json"
	"net/http"
)

// Respond convert um valor Go em JSON e responde a requisição ao cliente
func Respond(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	// armazena o status code da requisição no contexto para que possa ser utilizado pelos middlewares
	SetStatusCode(ctx, statusCode)

	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
