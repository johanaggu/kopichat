package mysql

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/johanaggu/kopichat/internal/chat"
)

type Conf struct {
	User string
	Pass string
	Host string
	Port string
	Name string
}
type Conn interface {
	RunMigrations(schema string) error
	Close() error
	RetrieveChat(ctx context.Context, id string) (*chat.Chat, error)
	CreateChat(ctx context.Context, ch *chat.Chat)(int64, error) 
}

func NewConn(ctx context.Context, c Conf) (Conn, error) {
	if c.User == "" || c.Pass == "" || c.Host == "" || c.Port == "" || c.Name == "" {
		return nil, fmt.Errorf("missing mysql config")
	}

	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.User, c.Pass, c.Host, c.Port, c.Name)

	db, err := sqlx.Connect("mysql", uri)
	if err != nil {
		return nil, err

	}

	return &conn{db}, nil
}

type conn struct {
	*sqlx.DB
}

func (c *conn) RunMigrations(schema string) error {
	c.MustExec(schema)
	return nil
}

func (c *conn) Close() error {
	return c.DB.Close()
}

func (c *conn) RetrieveChat(ctx context.Context, id string) (*chat.Chat, error) {
	var dbChat DBChat
	err := c.GetContext(ctx, &dbChat, "SELECT chat FROM chats WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return chat.NewChat(dbChat.ID, dbChat.APIID), nil
}

func (c *conn) CreateChat(ctx context.Context, ch *chat.Chat) (int64, error) {
	tx := c.MustBegin()
	result := tx.MustExec("INSERT INTO chats (api_id) VALUES ($1)", ch.ExternalID())
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return result.LastInsertId()
}
