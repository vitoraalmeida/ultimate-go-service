package summary

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
	UserID   *uuid.UUID `validate:"omitempty,uuid4"`
	UserName *string    `validate:"omitempty,min=3"`
}

// Validate checks the data in the model is considered clean.
// Validate checa se o dado está no formato correto
func (qf *QueryFilter) Validate() error {
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// WithUserID define o campo ID para ser usado no filtro
func (qf *QueryFilter) WithUserID(userID uuid.UUID) {
	qf.UserID = &userID
}

// WithUserName define o campo UserName para ser usado no filtro
func (qf *QueryFilter) WithUserName(userName string) {
	qf.UserName = &userName
}
