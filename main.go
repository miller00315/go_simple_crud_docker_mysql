package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" //Importação implicita que irá fornecer o driver para o mysql
)

func main() {
	//usuário:senha@/banco de dados?configurações adicionais
	connectionString := "root:RootPassword@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, error := sql.Open("mysql", connectionString)

	if error != nil {
		log.Fatal(error)
	}

	defer db.Close() // Fecha o banco de dados antes da função main terminar

	// Verifica se a conexão com o banco de dados ocorreu corretamente
	// Reaproveito a variavel error utilizando apenas =
	if error = db.Ping(); error != nil {
		log.Fatal(error)
	}

	fmt.Println("Conectado ao banco uhuul")

	lines, error := db.Query("select * from users")

	if error != nil {
		log.Fatal(error)
	}

	defer lines.Close() // Libera o cursor das linhas

	fmt.Println(lines)
}
