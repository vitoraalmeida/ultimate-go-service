// Package summary implementa um pacote business que implementa uma view.
// View aqui serve para unir dados que estão presentes em 2 domínios diferentes
package summary

import (
	"context"
	"fmt"

	"github.com/vitoraalmeida/service/business/data/order"
)

type Storer interface {
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Summary, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}

// =============================================================================

type Core struct {
	storer Storer
}

func NewCore(storer Storer) *Core {
	return &Core{
		storer: storer,
	}
}

func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Summary, error) {
	users, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return users, nil
}

func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}
