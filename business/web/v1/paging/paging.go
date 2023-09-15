// Package paging provê paginação para resições
package paging

import (
	"net/http"
	"strconv"

	"github.com/vitoraalmeida/service/business/sys/validate"
)

// Response é o retorno quando um query é executada
type Response[T any] struct {
	Items       []T `json:"items"`
	Total       int `json:"total"`
	Page        int `json:"page"`
	RowsPerPage int `json:"rowsPerPage"`
}

// NewResponse constrói uma resposta paginada
func NewResponse[T any](items []T, total int, page int, rowsPrePage int) Response[T] {
	return Response[T]{
		Items:       items,
		Total:       total,
		Page:        page,
		RowsPerPage: rowsPrePage,
	}
}

// =============================================================================

// Page representa a página e e linhas da paǵina
type Page struct {
	Number      int
	RowsPerPage int
}

// ParseRequest faz o parse do request recuperando as informações passadas
// relativas à paginação desejada pelo usuário
func ParseRequest(r *http.Request) (Page, error) {
	values := r.URL.Query()

	number := 1
	if page := values.Get("page"); page != "" {
		var err error
		number, err = strconv.Atoi(page)
		if err != nil {
			return Page{}, validate.NewFieldsError("page", err)
		}
	}

	rowsPerPage := 10
	if rows := values.Get("rows"); rows != "" {
		var err error
		rowsPerPage, err = strconv.Atoi(rows)
		if err != nil {
			return Page{}, validate.NewFieldsError("rows", err)
		}
	}

	return Page{
		Number:      number,
		RowsPerPage: rowsPerPage,
	}, nil
}
