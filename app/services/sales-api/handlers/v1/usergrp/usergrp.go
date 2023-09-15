// Package usergrp maintains the group of handlers for user access.
package usergrp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/vitoraalmeida/service/business/core/user"
	v1 "github.com/vitoraalmeida/service/business/web/v1"
	"github.com/vitoraalmeida/service/business/web/v1/paging"
	"github.com/vitoraalmeida/service/foundation/web"
)

// Handlers manages the set of user endpoints.
type Handlers struct {
	user *user.Core
}

// New constructs a handlers for route access.
func New(user *user.Core) *Handlers {
	return &Handlers{
		user: user,
	}
}

// Create adiciona um novo usuário no sistema
func (h *Handlers) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// cria uma instância do objeto para criação de usuários
	var app AppNewUser
	// faz a conversão do JSON para o objeto e chama a validação internamente
	if err := web.Decode(r, &app); err != nil {
		return err
	}

	// converte para o modelo de criação de novo usuário interno do domínio
	nc, err := toCoreNewUser(app)
	if err != nil {
		return v1.NewRequestError(err, http.StatusBadRequest)
	}

	usr, err := h.user.Create(ctx, nc)
	if err != nil {
		if errors.Is(err, user.ErrUniqueEmail) {
			// erro confiável, conhecido, que retorna mensagem mais específica para o usuário
			return v1.NewRequestError(err, http.StatusConflict)
		}
		// erros que o usuário final não deve ter detalhes
		return fmt.Errorf("create: usr[%+v]: %w", usr, err)
	}

	return web.Respond(ctx, w, toAppUser(usr), http.StatusCreated)
}

// considera que o id está vindo do token JWT
// Update updates a user in the system.
// func (h *Handlers) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
// 	var app AppUpdateUser
// 	if err := web.Decode(r, &app); err != nil {
// 		return err
// 	}

// 	userID := auth.GetUserID(ctx)

// 	usr, err := h.user.QueryByID(ctx, userID)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, user.ErrNotFound):
// 			return v1.NewRequestError(err, http.StatusNotFound)
// 		default:
// 			return fmt.Errorf("querybyid: userID[%s]: %w", userID, err)
// 		}
// 	}

// 	uu, err := toCoreUpdateUser(app)
// 	if err != nil {
// 		return v1.NewRequestError(err, http.StatusBadRequest)
// 	}

// 	usr, err = h.user.Update(ctx, usr, uu)
// 	if err != nil {
// 		return fmt.Errorf("update: userID[%s] uu[%+v]: %w", userID, uu, err)
// 	}

// 	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
// }

// Delete removes a user from the system.
// func (h *Handlers) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
// 	userID := auth.GetUserID(ctx)

// 	usr, err := h.user.QueryByID(ctx, userID)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, user.ErrNotFound):
// 			return web.Respond(ctx, w, nil, http.StatusNoContent)
// 		default:
// 			return fmt.Errorf("querybyid: userID[%s]: %w", userID, err)
// 		}
// 	}

// 	if err := h.user.Delete(ctx, usr); err != nil {
// 		return fmt.Errorf("delete: userID[%s]: %w", userID, err)
// 	}

// 	return web.Respond(ctx, w, nil, http.StatusNoContent)
// }

// Query retorna uma lista de usuários paginada
func (h *Handlers) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	// Faz o parsing por informações de paginção
	page, err := paging.ParseRequest(r)
	if err != nil {
		return err
	}

	// Faz o parsing por informações de filtragem de resultados
	// filter é um objeto com informações de usuário que podem
	// ser usadas para filtrar resultados
	filter, err := parseFilter(r)
	if err != nil {
		return err
	}

	// Faz o parsing por informações de ordenação de resultados
	orderBy, err := parseOrder(r)
	if err != nil {
		return err
	}

	// executa a query no banco com base nas informações passadas
	users, err := h.user.Query(ctx, filter, orderBy, page.Number, page.RowsPerPage)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	items := make([]AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	total, err := h.user.Count(ctx, filter)
	if err != nil {
		return fmt.Errorf("count: %w", err)
	}

	return web.Respond(ctx, w, paging.NewResponse(items, total, page.Number, page.RowsPerPage), http.StatusOK)
}

// considera que o id está vindo do token JWT
// QueryByID returns a user by its ID.
// func (h *Handlers) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
// 	id := auth.GetUserID(ctx)

// 	usr, err := h.user.QueryByID(ctx, id)
// 	if err != nil {
// 		switch {
// 		case errors.Is(err, user.ErrNotFound):
// 			return v1.NewRequestError(err, http.StatusNotFound)
// 		default:
// 			return fmt.Errorf("querybyid: id[%s]: %w", id, err)
// 		}
// 	}

// 	return web.Respond(ctx, w, toAppUser(usr), http.StatusOK)
// }
