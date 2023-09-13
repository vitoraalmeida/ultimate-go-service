package user

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/vitoraalmeida/service/business/sys/validate"
)

// QueryFilter agrupa os campos disponíveis pelos quais uma consulta pode ser filtrada
type QueryFilter struct {
	// utiliza ponteiros para dar a possibilidade de deixar um ou mais campos vazios (nil)
	// e podermos passar o objeto inteiro para que seja utilizado com base nos campos
	// que não forem nulos
	ID               *uuid.UUID    `validate:"omitempty,uuid4"` // tags para validação do pacote validator go-playground/validator/v10"
	Name             *string       `validate:"omitempty,min=3"`
	Email            *mail.Address `validate:"omitempty,email"`
	StartCreatedDate *time.Time    `validate:"omitempty"`
	EndCreatedDate   *time.Time    `validate:"omitempty"`
}

// Validate checa se o dado está limpo
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// WithUserID define o campo ID para ser usado no filtro
func (qf *QueryFilter) WithUserID(userID uuid.UUID) {
	qf.ID = &userID
}

// WithName define o campo Name para ser usado no filtro
func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}

// WithEmail define o campo Email para ser usado no filtro
func (qf *QueryFilter) WithEmail(email mail.Address) {
	qf.Email = &email
}

// WithStartDateCreated define o campo StartDateCreated para ser usado no filtro
func (qf *QueryFilter) WithStartDateCreated(startDate time.Time) {
	d := startDate.UTC()
	qf.StartCreatedDate = &d
}

// WithEndCreatedDate define o campo EndCreatedDate para ser usado no filtro
func (qf *QueryFilter) WithEndCreatedDate(endDate time.Time) {
	d := endDate.UTC()
	qf.EndCreatedDate = &d
}
