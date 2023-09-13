package summary

import "github.com/google/uuid"

// Summary representa informação sobre um usuário e seus produtos
// Foi gerado para que pudéssemos atender a necessidade fazer o
// relatório de produtos em conjunto com os nomes dos usuário que
// os registraram
type Summary struct {
	UserID     uuid.UUID
	UserName   string  // vem de user
	TotalCount int     // vem de products
	TotalCost  float64 // vem de Products
}
