package userdb

import (
	"fmt"

	"github.com/vitoraalmeida/service/business/core/user"
	"github.com/vitoraalmeida/service/business/data/order"
)

var orderByFields = map[string]string{
	user.OrderByID:      "user_id",
	user.OrderByName:    "name",
	user.OrderByEmail:   "email",
	user.OrderByRoles:   "roles",
	user.OrderByEnabled: "enabled",
}

// adiciona na query que vai ser executada a parte da ordenação
// caso seja necessário
func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
