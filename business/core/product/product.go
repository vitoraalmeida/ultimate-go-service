// Package product provides an example of a core business API. Right now these
// calls are just wrapping the data/store layer. But at some point you will
// want auditing or something that isn't specific to the data/store layer.
package product

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vitoraalmeida/service/business/core/user"
	"github.com/vitoraalmeida/service/business/data/order"
	"go.uber.org/zap"
)

// Conjunto de erros para operações CRUD
var (
	ErrNotFound = errors.New("product not found")
)

// Abstrai qual é a implementação de fato que vai gerenciar a interção
// com o armazenamento de usuário, desde que possua esse comportamento
type Storer interface {
	Create(ctx context.Context, prd Product) error
	Update(ctx context.Context, prd Product) error
	Delete(ctx context.Context, prd Product) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (Product, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error)
}

// Core é a API para o domínio Product, gerencia as ações num produto
type Core struct {
	// Abstrai qual é a implementação de fato que vai gerenciar a interção
	// com o armazenamento de usuário
	log     *zap.SugaredLogger
	usrCore *user.Core
	storer  Storer
}

// NewCore constrói Core para uso da API de produtos
func NewCore(log *zap.SugaredLogger, usrCore *user.Core, storer Storer) *Core {
	core := Core{
		log:     log,
		usrCore: usrCore,
		storer:  storer,
	}

	return &core
}

// Create insere um novo produto no banco de dados, retornando o produto com o ID que foi gerado pelo sistema
// semantica de ponteiro para APIs         semantica de valor para Dados e para interfaces (context.Context)
func (c *Core) Create(ctx context.Context, np NewProduct) (Product, error) {
	now := time.Now()

	prd := Product{
		ID:          uuid.New(),
		Name:        np.Name,
		Cost:        np.Cost,
		Quantity:    np.Quantity,
		UserID:      np.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := c.storer.Create(ctx, prd); err != nil {
		return Product{}, fmt.Errorf("create: %w", err)
	}

	return prd, nil
}

// invalid or does not reference an existing Product.
// Update modiffica dados sobre um produto. Retorna erro se o ID especificado não existir ou não for um uuid valido
func (c *Core) Update(ctx context.Context, prd Product, up UpdateProduct) (Product, error) {
	if up.Name != nil {
		prd.Name = *up.Name
	}
	if up.Cost != nil {
		prd.Cost = *up.Cost
	}
	if up.Quantity != nil {
		prd.Quantity = *up.Quantity
	}
	prd.DateUpdated = time.Now()

	if err := c.storer.Update(ctx, prd); err != nil {
		return Product{}, fmt.Errorf("update: %w", err)
	}

	return prd, nil
}

func (c *Core) Delete(ctx context.Context, prd Product) error {
	if err := c.storer.Delete(ctx, prd); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query busca todos os produtos do banco com paginação
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]Product, error) {
	prds, err := c.storer.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}

// Count retorna o numero total de produtos no banco
func (c *Core) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return c.storer.Count(ctx, filter)
}

func (c *Core) QueryByID(ctx context.Context, productID uuid.UUID) (Product, error) {
	prd, err := c.storer.QueryByID(ctx, productID)
	if err != nil {
		return Product{}, fmt.Errorf("query: productID[%s]: %w", productID, err)
	}

	return prd, nil
}

// QueryByUserID busca um produto que foi adicionado por um usuário determinado
func (c *Core) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error) {
	prds, err := c.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return prds, nil
}
