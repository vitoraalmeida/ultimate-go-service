// Package order provides support for describing the ordering of data.
package order

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/vitoraalmeida/service/business/sys/validate"
)

// Define direções para a ordenação
const (
	ASC  = "ASC"
	DESC = "DESC"
)

var directions = map[string]string{
	ASC:  "ASC",
	DESC: "DESC",
}

// =============================================================================

// By representa um campo para ser ordenado e uma direção
type By struct {
	Field     string
	Direction string
}

// NewBy contrói um By
func NewBy(field string, direction string) By {
	return By{
		Field:     field,
		Direction: direction,
	}
}

// =============================================================================

// Parse constrói um valor By fazendo um parsing de uma string na forma "field,direction"
func Parse(r *http.Request, defaultOrder By) (By, error) {
	// busca infos de como o usuário solicitou a filtragem/ordenação
	v := r.URL.Query().Get("orderBy")

	if v == "" {
		return defaultOrder, nil
	}

	orderParts := strings.Split(v, ",")

	var by By
	switch len(orderParts) {
	case 1:
		// se não foi passado uma direção, usamos ASC por padrão
		by = NewBy(strings.Trim(orderParts[0], " "), ASC)
	case 2:
		by = NewBy(strings.Trim(orderParts[0], " "), strings.Trim(orderParts[1], " "))
	default:
		return By{}, validate.NewFieldsError(v, errors.New("unknown order field"))
	}

	if _, exists := directions[by.Direction]; !exists {
		return By{}, validate.NewFieldsError(v, fmt.Errorf("unknown direction: %s", by.Direction))
	}

	return by, nil
}
