package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	db "github.com/johanaggu/kopichat/internal/infrastructure/mysql"
)

var schema = `
	CREATE TABLE IF NOT EXISTS chats(   
		id    INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,   
		api_id VARCHAR(255) NOT NULL 
	);
`

type Message struct {
	Message string `json:"message"`
}

func Migrate(ctx context.Context) (Message, error) {
	root, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&charset=utf8mb4,utf8",
			os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT")))
	if err != nil {
		return Message{
			Message: "ailed to create db connection",
		}, fmt.Errorf("failed to create db connection")
	}
	root.ExecContext(ctx, `CREATE DATABASE IF NOT EXISTS kopichat`)

	conn, err := db.NewConn(ctx, db.Conf{
		User: os.Getenv("MYSQL_USER"),
		Pass: os.Getenv("MYSQL_PASSWORD"),
		Host: os.Getenv("MYSQL_HOST"),
		Port: os.Getenv("MYSQL_PORT"),
		Name: os.Getenv("MYSQL_DATABASE"),
	})
	if err != nil {
		return Message{
			Message: "failed to create db connection",
		}, fmt.Errorf("failed to create db connection")
	}

	defer conn.Close()

	if err := conn.RunMigrations(schema); err != nil {
		return Message{
			Message: "ailed to run migrations ",
		}, fmt.Errorf("failed to run migrations ")
	}

	return Message{
		Message: "migrations run successfully",
	}, nil
}
