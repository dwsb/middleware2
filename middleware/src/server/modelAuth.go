package modelAuth

import (
	"crypto/rand"
	"crypto/sha1"
	"fmt"
)

type User struct {
	email string
	senha string
}

type DbUser struct {
	meuSlice []User
}

type Session struct {
	user  User
	token string
	ip    string
}

type Dbsession struct {
	sessions []Session
}

func hash(b []byte) string {
	h := sha1.New()
	h.Write(b)
	sum := h.Sum(nil)
	armored := fmt.Sprintf("%x", sum)
	return armored
}

//gera o token
func randomString() (string, error) {
	b := make([]byte, 10)
	_, err := rand.Read(b)

	if err != nil {
		fmt.Printf("erro: %v", err)
		return "", err
	}

	armored := hash(b)
	return armored, err
}

//funcao para verificar login
func loginSessionUser(db *DbUser, dbSession *Dbsession, user User) Session {
	for i := 0; i < (len((*db).meuSlice)); i++ {
		if (*db).meuSlice[i].email == user.email && (*db).meuSlice[i].senha == user.senha {
			fmt.Println("Fez login!!!")
			return addSession(dbSession, user)

		}
	}
	fmt.Println("Login falhou!!!")
	return Session{}
}

//funcao para verificar se o usuario ja existe
func noDuplicateUser(db *DbUser, email string) bool {
	for i := 0; i < (len((*db).meuSlice)); i++ {
		if (*db).meuSlice[i].email == email {
			return false
		}
	}
	return true
}

//funcao para nao duplicar a sessao de um usuario ja existente
func noDuplicateSession(db *Dbsession, email string, token string) bool {
	for i := 0; i < (len((*db).sessions)); i++ {
		if (*db).sessions[i].user.email == email || (*db).sessions[i].token == token {
			return false
		}
	}
	return true
}

//funcao pra remover o usuario
func rmUser(db *DbUser, user User) bool {
	for i := 0; i < (len((*db).meuSlice)); i++ {
		if (*db).meuSlice[i].email == user.email && (*db).meuSlice[i].senha == user.senha {
			copy((*db).meuSlice[i:], (*db).meuSlice[i+1:])
			(*db).meuSlice[len((*db).meuSlice)-1] = User{}
			(*db).meuSlice = (*db).meuSlice[:len((*db).meuSlice)-1]
			return true
		}
	}
	return false
}

//adiciona usuario
func addUser(db *DbUser, user User) bool {
	if noDuplicateUser(db, user.email) {
		(*db).meuSlice = append((*db).meuSlice, user)
		fmt.Println("Adicionado com sucesso o usuario:", user.email)
		return true
	}
	fmt.Println("Não conseguiu adicionar o usuario:", user.email)
	return false
}

//update de senha
func updatePasswUser(db *DbUser, user User) bool {
	for i := 0; i < (len((*db).meuSlice)); i++ {
		if (*db).meuSlice[i].email == user.email {
			(*db).meuSlice[i].senha = user.senha
			return true
		}
	}
	return false
}

//adicionar sessao
func addSession(db *Dbsession, user1 User) Session {
	token1, _ := randomString()
	for !noDuplicateSession(db, user1.email, token1) {
		token1, _ = randomString()
	}
	aux := Session{user: user1, token: token1, ip: ""}
	(*db).sessions = append((*db).sessions, aux)
	fmt.Println("sessao adicionada")
	return aux
}

//fazer logout da sessao
func logoutSessionUser(db *Dbsession, user1 User, token1 string) bool {
	for i := 0; i < len((*db).sessions); i++ {
		if user1.email == (*db).sessions[i].user.email && (*db).sessions[i].token == token1 {
			copy((*db).sessions[i:], (*db).sessions[i+1:])
			(*db).sessions[len((*db).sessions)-1] = Session{}
			(*db).sessions = (*db).sessions[:len((*db).sessions)-1]
			fmt.Println("Fazendo logout:", user1.email)
			return true
		}
	}
	return false
}
