package testgrp

import (
	"context"
	"net/http"

	"github.com/vitoraalmeida/service/foundation/web"
)

// Rota de exemplo

func Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// handler devem
	// validar os dados
	// processar na camada de business
	// lidar com erros

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	// lidar com a resposta OK
	return web.Respond(ctx, w, status, http.StatusOK)
}
