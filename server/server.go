package server

import (
	"database/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type user struct {
	ID    uint32 `json:id`
	Name  string `json:name`
	Email string `json:email`
}

// Cria o usuário
func CreateUser(w http.ResponseWriter, r *http.Request) {
	requestBody, error := ioutil.ReadAll(r.Body)

	if error != nil {
		w.Write([]byte("Falha ao ler o corpo da requisição"))
		return
	}

	var user user

	if error = json.Unmarshal(requestBody, &user); error != nil {
		w.Write([]byte("Falha ao converter o user"))
		return
	}

	db, error := database.Connect()

	if error != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados"))
		return
	}

	defer db.Close() // Fecha o banco de dados ao terminar tudo, importante!!!

	statement, error := db.Prepare("insert into users (name, email) values (?, ?)")

	if error != nil {
		w.Write([]byte("Erro ao criar o statement"))
		return
	}

	defer statement.Close() // Fecha o statement ao terminar tudo, importante !!!

	insertion, error := statement.Exec(user.Name, user.Email) // Executa o statement para salvar o usuário

	if error != nil {
		w.Write([]byte("Erro ao executar o statement"))
		return
	}

	insertionId, error := insertion.LastInsertId() // Recupera o ultimo id inserido

	if error != nil {
		w.Write([]byte("Erro ao recuperar id"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("Success insert user, Id: %d", insertionId)))
}
