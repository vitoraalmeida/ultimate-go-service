package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	"github.com/vitoraalmeida/service/foundation/web"
)

// Rota de exemplo

func Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// handler devem
	// validar os dados
	// processar na camada de business

	// lidar com erros
	// gera um erro fictício para testar o error handling
	if n := rand.Intn(100); n%2 == 0 {
		return errors.New("UNTRUSTED ERROR")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	// lidar com a resposta OK
	return web.Respond(ctx, w, status, http.StatusOK)
}
