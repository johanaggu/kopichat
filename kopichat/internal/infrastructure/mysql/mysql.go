package mysql

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/johanaggu/kopichat/internal/model/chat"
	"github.com/johanaggu/kopichat/internal/model/message"
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

	SaveChat(ctx context.Context, externalID string) (*chat.Chat, error)
	GetChatByID(ctx context.Context, id string) (*chat.Chat, error)

	SaveMessage(ctx context.Context, role message.Role, content, chatID string) (*message.Message, error)
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

func (c *conn) GetChatByID(ctx context.Context, id string) (*chat.Chat, error) {
	var dbChat DBChat
	err := c.GetContext(ctx, &dbChat, "SELECT id, external_id FROM chats WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	return chat.New(dbChat.ID, dbChat.ExternalID), nil
}

func (c *conn) SaveChat(ctx context.Context, externalID string) (*chat.Chat, error) {
	tx := c.MustBegin()

	var id string
	if err := tx.GetContext(ctx, &id, `SELECT UUID()`); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("select uuid: %w", err)
	}

	result := tx.QueryRowContext(ctx, "INSERT INTO chats (id, external_id) VALUES (?, ?)", id, externalID)
	if result.Err() != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating chat: %w", result.Err())
	}

	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return chat.New(id, externalID), nil
}

func (c *conn) SaveMessage(ctx context.Context, role message.Role, content, chatID string) (*message.Message, error) {
	tx := c.MustBegin()
	result := tx.MustExecContext(ctx, `INSERT INTO messages 
		(role, content, chat_id) VALUES (?, ?, ?)`, string(role), content, chatID)
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error retrieve last inserted message id")
	}

	return message.NewMessage(int(id), role, content, chatID), nil
}
