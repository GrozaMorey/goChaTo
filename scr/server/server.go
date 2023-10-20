package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type userBody struct {
	Usernama string `json:"username"`
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
		return
	}
	_, err = srv.Pool.Exec(context.Background(), "INSERT INTO users (username, pass) VALUES ($1, $2)", userBody.Usernama, userBody.Password)
	if err != nil {
		log.Println(err)
	}
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
