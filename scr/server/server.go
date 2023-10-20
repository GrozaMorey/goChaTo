package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type userBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Server struct {
	Pool *pgxpool.Pool
}

func New() (*Server, error) {
	var srv Server
	dbpool, err := pgxpool.New(context.Background(), "postgres://postgres:123@localhost:5432/Go")
	if err != nil {
		log.Fatal("cant connect to db")
	}
	srv.Pool = dbpool
	return &srv, err
}

func (srv *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var userBody userBody
	err := json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		log.Println("body is empty", err)
		http.Error(w, "body is empty", http.StatusBadRequest)
		return
	}

	if userBody.Password != "" && userBody.Username != "" {

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(userBody.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Something wrong w/ hashing password")
		}

		_, err = srv.Pool.Exec(context.Background(), "INSERT INTO users (username, pass) VALUES ($1, $2)", userBody.Username, hashPassword)
		if err != nil {
			log.Println(err)
		}
		w.Write([]byte("Success registarion"))
		return
	}
	http.Error(w, "Username and Password required", http.StatusBadRequest)
}

func (srv *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type Row struct {
		Id       int
		Username string
		Pass     string
	}
	var userBody userBody

	err := json.NewDecoder(r.Body).Decode(&userBody)
	if err != nil {
		log.Println("body is empty", err)
		http.Error(w, "body is empty", http.StatusBadRequest)
		return
	}
	if userBody.Password != "" && userBody.Username != "" {
		rows, err := srv.Pool.Query(context.Background(), "select * from users where username = $1", userBody.Username)
		if err != nil {
			log.Println(err)
			return
		}
		user, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Row])
		if err != nil {
			log.Println("Cant parse raw", err)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(userBody.Password))
		if err != nil {
			log.Println(err)
			http.Error(w, "Password is incorrect", http.StatusBadRequest)
			return
		}

		w.Write([]byte("Success login"))
		return
	}
	http.Error(w, "Username and Password required", http.StatusBadRequest)
}

// func main() {
// 	dbpool, err := New()
// 	_, err = dbpool.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users(
// 		ID serial primary key,
// 		username text,
// 		pass text
// 	); CREATE TABLE IF NOT EXISTS chats(
// 		ID serial primary key
// 	); CREATE TABLE IF NOT EXISTS messages(
// 		ID serial primary key,
// 		from_id integer references users(ID),
// 		chat_id integer references chats(ID),
// 		content text,
// 		timestamp timestamp
// 	); CREATE TABLE IF NOT EXISTS user_chat(
// 		user_id integer references users(ID),
// 		chat_id integer references chats(ID)
// 	)`)
// 	if err != nil {
// 		log.Print(err)
// 	}

// }
