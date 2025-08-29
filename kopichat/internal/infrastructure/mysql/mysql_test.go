package mysql

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveChat(t *testing.T) {
	ctx := context.Background()

	cfg := Conf{
		User: os.Getenv("MYSQL_USER"),
		Pass: os.Getenv("MYSQL_PASSWORD"),
		Host: os.Getenv("MYSQL_HOST"),
		Port: os.Getenv("MYSQL_PORT"),
		Name: os.Getenv("MYSQL_DATABASE"),
	}

	connection, err := NewConn(t.Context(), cfg)
	if err != nil {
		t.Fatal(err)
	}

	externalID := "external_id_1"
	actualChat, err := connection.SaveChat(ctx, externalID)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEmpty(t, actualChat.ID())
	assert.Equal(t, externalID, actualChat.ExternalID())
}
