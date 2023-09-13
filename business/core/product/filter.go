package product

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vitoraalmeida/service/business/sys/validate"
)

// QueryFilter agrupa os campos disponíveis pelos quais uma consulta pode ser filtrada. É passado por query params numa url
type QueryFilter struct {
	// utiliza ponteiros para dar a possibilidade de deixar um ou mais campos vazios (nil)
	// e podermos passar o objeto inteiro para que seja utilizado com base nos campos
	// que não forem nulos
	ID       *uuid.UUID `validate:"omitempty"`
	Name     *string    `validate:"omitempty,min=3"`
	Cost     *float64   `validate:"omitempty,numeric"`
	Quantity *int       `validate:"omitempty,numeric"`
}

// Validate checa se o dado está no formato correto
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// WithProductID define o campo ID para ser usado no filtro
func (qf *QueryFilter) WithProductID(productID uuid.UUID) {
	qf.ID = &productID
}

// WithName define o campo Name para ser usado no filtro
func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}

// WithCost define o campo Cost para ser usado no filtro
func (qf *QueryFilter) WithCost(cost float64) {
	qf.Cost = &cost
}

// WithQuantity define o campo Quantity para ser usado no filtro
func (qf *QueryFilter) WithQuantity(quantity int) {
	qf.Quantity = &quantity
}
