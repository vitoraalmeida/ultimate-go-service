package product

import (
	"time"

	"github.com/google/uuid"
)

// Product representa um produto individual.
type Product struct {
	ID          uuid.UUID
	Name        string
	Cost        float64
	Quantity    int
	Sold        int
	Revenue     int
	UserID      uuid.UUID // indica uma relação com User. O usuário que registrou esse produto
	DateCreated time.Time
	DateUpdated time.Time
}

// NewProduct representa o modelo de dados que exigimos do cliente para criar um produto
type NewProduct struct {
	Name     string
	Cost     float64
	Quantity int
	UserID   uuid.UUID
}

// UpdateUser  contém informação necessária para atualizar dados de um usuário
type UpdateProduct struct {
	Name     *string  // Usamos a semântica de ponteiro aqui
	Cost     *float64 // para mostrar que alguns desses dados
	Quantity *int     // podem ser nulos na tentativa de atualizar
} // um usuário, pois podemos querer atualizar
// apenas algum dos dados do usuário
