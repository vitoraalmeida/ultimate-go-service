package usergrp

import (
	"errors"
	"net/http"

	"github.com/vitoraalmeida/service/business/core/user"
	"github.com/vitoraalmeida/service/business/cview/user/summary"
	"github.com/vitoraalmeida/service/business/data/order"
	"github.com/vitoraalmeida/service/business/sys/validate"
)

// conjunto de todos os campos possíveis pelos quais podemos ordenar os resultados
var orderByFields = map[string]struct{}{
	user.OrderByID:      {},
	user.OrderByName:    {},
	user.OrderByEmail:   {},
	user.OrderByRoles:   {},
	user.OrderByEnabled: {},
}

func parseOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, user.DefaultOrderBy)
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	return orderBy, nil
}

// =============================================================================

// conjunto de todos os campos possíveis pelos quais podemos ordenar os resultados
var orderBySummaryFields = map[string]struct{}{
	summary.OrderByUserID:   {},
	summary.OrderByUserName: {},
}

func parseSummaryOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, user.DefaultOrderBy)
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderBySummaryFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	return orderBy, nil
}
