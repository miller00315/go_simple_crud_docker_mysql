package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Importação implicita que irá fornecer o driver para o mysql
)

// Connect conecta ao banco de dados
func Connect() (*sql.DB, error) {

	// usuário:senha@/banco de dados?configurações adicionais
	connectionString := "root:RootPassword@/devbook?charset=utf8&parseTime=True&loc=Local"

	// Recebe um banco ou um erro
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	// Verifica se a conexão com o banco de dados ocorreu corretamente
	// Reaproveito a variavel error utilizando apenas =
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
