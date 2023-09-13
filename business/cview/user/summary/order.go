package summary

import "github.com/vitoraalmeida/service/business/data/order"

var DefaultOrderBy = order.NewBy(OrderByUserID, order.ASC)

const (
	OrderByUserID   = "userid"
	OrderByUserName = "userName"
)
