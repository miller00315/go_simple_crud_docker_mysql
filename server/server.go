package server

import (
	"database/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type user struct {
	ID    uint32 `json:id`
	Name  string `json:name`
	Email string `json:email`
}

// Cria o usuário
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Coleta o corpo da requisição
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

// Retorna todos os usuaŕios
func GetUsers(w http.ResponseWriter, r *http.Request) {
	db, error := database.Connect()

	if error != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados"))
		return
	}

	defer db.Close() // Fecha o banco de dados ao terminar tudo, importante!!!

	lines, error := db.Query("select * from users")

	if error != nil {
		w.Write([]byte("Falha ao buscar usuarios"))
		return
	}

	defer lines.Close() // Fecha as linhas da query ao finalizar tudo, importante

	var users []user // Cria um slice de usuarios

	for lines.Next() {
		var user user // Cria um único usuário

		if error := lines.Scan(&user.ID, &user.Name, &user.Email); error != nil {
			w.Write([]byte("Erro ao Escanear usuários"))
			return
		}

		users = append(users, user) // O slice é populado
	}

	w.WriteHeader(http.StatusOK) // Escreve o cabeçalho da resposta

	// Converte para json e envia a resposta
	if error := json.NewEncoder(w).Encode(users); error != nil {
		w.Write([]byte("Erro ao converter para json"))
		return
	}

}

// Retorna um usuário
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Pega o id recebido na requisição
	ID, error := strconv.ParseUint(params["id"], 10, 32)

	if error != nil {
		fmt.Println(error)
		w.Write([]byte("Erro ao converter para numero para inteiro"))
		return
	}

	// Abro a conexão depois de verificar que o id está correto
	db, error := database.Connect()

	if error != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados"))
		return
	}

	defer db.Close() // Fecha o banco de dados ao terminar tudo, importante!!!

	// Coleta o usuário de uma linha especifica
	line, error := db.Query("select * from users where id = ?", ID)

	if error != nil {
		w.Write([]byte("Falha ao buscar usuario"))
		return
	}

	defer line.Close() // Fecha a linha ao final da execução, importante!!!

	var user user // Cria um único usuário

	if line.Next() {
		// Escaneia o usuário
		if error := line.Scan(&user.ID, &user.Name, &user.Email); error != nil {
			w.Write([]byte("Erro ao Escanear usuário"))
			return
		}
	}

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Usuário não encontrado"))
		return
	}

	w.WriteHeader(http.StatusOK) // Escreve o cabeçalho da resposta

	// Converte para json e envia a resposta
	if error := json.NewEncoder(w).Encode(user); error != nil {
		w.Write([]byte("Erro ao converter para json"))
		return
	}
}

// Atualiza um usuário
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Mux é utilizado para coletar os parâmetros da reuqisição
	params := mux.Vars(r)

	// Pega o id recebido na requisição
	ID, error := strconv.ParseUint(params["id"], 10, 32)

	if error != nil {
		fmt.Println(error)
		w.Write([]byte("Erro ao converter para numero para inteiro"))
		return
	}

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

	// Abro a conexão depois de verificar que o id está correto
	db, error := database.Connect()

	if error != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados"))
		return
	}

	defer db.Close() // Fecha o banco de dados ao terminar tudo, importante!!!

	statement, error := db.Prepare("update users set name = ?, email = ? where id = ?")

	if error != nil {
		fmt.Println(error)
		w.Write([]byte("Falha criar statement"))
		return
	}

	defer statement.Close() // Fecha o statement ao final, importante

	// Executa o statement para salvar o usuário e já verifica se há erros
	// Deve utilizar o ID recebido na requisição
	if _, error := statement.Exec(user.Name, user.Email, ID); error != nil {
		w.Write([]byte("Erro ao atualizar o usuário"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Apaga um usuário do banco de dados
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Mux é utilizado para coletar os parâmetros da reuqisição
	params := mux.Vars(r)

	// Pega o id recebido na requisição
	ID, error := strconv.ParseUint(params["id"], 10, 32)

	if error != nil {
		fmt.Println(error)
		w.Write([]byte("Erro ao converter para numero para inteiro"))
		return
	}

	// Abro a conexão depois de verificar que o id está correto
	db, error := database.Connect()

	if error != nil {
		w.Write([]byte("Falha ao conectar ao banco de dados"))
		return
	}

	defer db.Close() // Fecha o banco de dados ao terminar tudo, importante!!!

	statement, error := db.Prepare("delete from users where id = ?")

	if error != nil {
		fmt.Println(error)
		w.Write([]byte("Falha criar statement"))
		return
	}

	defer statement.Close() // Fecha o statement ao final, importante

	// Executa o statement para salvar o usuário e já verifica se há erros
	// Deve utilizar o ID recebido na requisição
	if _, error := statement.Exec(ID); error != nil {
		w.Write([]byte("Erro ao deletar o usuário"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
