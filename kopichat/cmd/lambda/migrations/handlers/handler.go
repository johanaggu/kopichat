package handlers

import (
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	db "github.com/johanaggu/kopichat/internal/infrastructure/mysql"
)

// Actions:
//  - "reset-db"         -> DROP/CREATE database (uses MYSQL_DATABASE)
//  - "create-chats"     -> create chats table
//  - "create-messages"  -> create messages table
//  - "create-index"     -> create index on messages(chat_id)
//  - "migrate"          -> create chats + messages + index

const resetDB = `
	DROP TABLE IF EXISTS chats;
	DROP TABLE IF EXISTS messages;
`

var chatsTable = `
	CREATE TABLE IF NOT EXISTS chats (
  		id     			CHAR(36) PRIMARY KEY DEFAULT (UUID()),
  		external_id  	VARCHAR(255) NULL,
  		created_at  	TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
`
var messageTable = `
	CREATE TABLE messages (
  		id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  		chat_id         CHAR(36) NOT NULL,
  		content         TEXT NOT NULL,
  		role            VARCHAR(20) NOT NULL,
  		created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  		CONSTRAINT fk_chat_messages FOREIGN KEY (chat_id)
  		REFERENCES chats(id)
  		ON DELETE CASCADE
	);
`

var messageChatIDIndex = `CREATE INDEX idx_chat_messages ON messages(chat_id);`

type Message struct {
	Message string `json:"message"`
}

func Migrate(ctx context.Context) (Message, error) {
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

	err = conn.RunMigrations(resetDB)
	if err != nil {
		return Message{}, fmt.Errorf("failed to reset db")
	}

	err = conn.RunMigrations(chatsTable)
	if err != nil {
		return Message{}, fmt.Errorf("failed to reset db")
	}

	err = conn.RunMigrations(messageTable)
	if err != nil {
		return Message{}, fmt.Errorf("failed to reset db")
	}

	err = conn.RunMigrations(messageChatIDIndex)
	if err != nil {
		return Message{}, fmt.Errorf("failed to reset db")
	}

	return Message{
		Message: "migrations run successfully",
	}, nil
}
