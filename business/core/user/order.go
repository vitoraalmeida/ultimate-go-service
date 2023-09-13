package user

import "github.com/vitoraalmeida/service/business/data/order"

// DefaultOrderBy representa a forma padrão de ordenação
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Conjunto de campos que podem ser usado para ordenar os resultados
const (
	OrderByID      = "userid"
	OrderByName    = "name"
	OrderByEmail   = "email"
	OrderByRoles   = "roles"
	OrderByEnabled = "enabled"
)
