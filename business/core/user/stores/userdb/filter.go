package userdb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vitoraalmeida/service/business/core/user"
)

// applyFilter vai criar a pacela da query SELECT após o WHERE para que possamos
// filtrar, selecionar apenas os campos que quisermos com base num objeto que vai
// ser passado em que campos nulos não serão selecionados na query
// filter = objeto do modelo contendo os campos que queremos recuperar do banco
/*
 data = é um map que armazena informações que serão passadas na query
        mapeia os campos da tabela que estão sendo consultados com os campos
        do objeto do modelo em questão, que será usado para apontar quais
        dados serão passados na query
		caso o objeto em questão tenha sido:

User {
	Name: Vitor,
	Id: nil,
	email: nil,
}

		a linha do sql que será executada será
		... WHERE name = :name", data['name']

		ou


		rows, err = sqlx.NamedQueryContext(ctx, db, query, data)

		em que NamedQueryContext vai subtituir os dados na query

		data também é usado para salvar a query que foi executada e logar em
		formado legível
*/

// buf = a query SQL que será gerada e modificada aqui
func (s *Store) applyFilter(filter user.QueryFilter, data map[string]interface{}, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["user_id"] = *filter.ID
		// wc representa o conjunto de clausulas where que serão usadas (where clauses)
		wc = append(wc, "user_id = :user_id") // :user_id é o formato usado pelo sqlx para proteger contra injection
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if filter.Email != nil {
		data["email"] = (*filter.Email).String()
		wc = append(wc, "email = :email")
	}

	if filter.StartCreatedDate != nil {
		data["start_date_created"] = *filter.StartCreatedDate
		wc = append(wc, "date_created >= :start_date_created")
	}

	if filter.EndCreatedDate != nil {
		data["end_date_created"] = *filter.EndCreatedDate
		wc = append(wc, "date_created <= :end_date_created")
	}

	if len(wc) > 0 {
		// adicionamos o WHERE na query base
		buf.WriteString(" WHERE ")
		// adicionamos cada um dos campos que vamos selecionar
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
