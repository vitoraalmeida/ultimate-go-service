package product

import "github.com/vitoraalmeida/service/business/data/order"

// DefaultOrderBy representa a forma padrão de ordenação
var DefaultOrderBy = order.NewBy(OrderByProdID, order.ASC)

// Conjunto de campos que podem ser usado para ordenar os resultados
const (
	OrderByProdID   = "productid"
	OrderByName     = "name"
	OrderByCost     = "cost"
	OrderByQuantity = "quantity"
	OrderBySold     = "sold"
	OrderByRevenue  = "revenue"
	OrderByUserID   = "userid"
)
