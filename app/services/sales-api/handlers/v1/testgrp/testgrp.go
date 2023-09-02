package testgrp

import (
	"context"
	"encoding/json"
	"net/http"
)

// Rota de exemplo

func Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// handler devem
	// validar os dados
	// processar na camada de business
	// retornar erros
	// lidar com a resposta OK
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return json.NewEncoder(w).Encode(status)
}
