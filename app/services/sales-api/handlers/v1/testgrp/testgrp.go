package testgrp

import (
	"encoding/json"
	"net/http"
)

// Rota de exemplo

func Test(w http.ResponseWriter, r *http.Request) {
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

	json.NewEncoder(w).Encode(status)
}
