package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConversation_Success(t *testing.T) {
	fakeAPIKey := "fake_api_key"
	conversationID := "conversation_id"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica método y path
		if r.Method != http.MethodPost {
			t.Fatal("unexpected http method")
		}
		if r.URL.Path != "/conversations" {
			t.Fatal("unexpected path")
		}

		wantAuth := "Bearer " + fakeAPIKey
		if auth := r.Header.Get("Authorization"); auth != wantAuth {
			t.Fatal("api key not found")
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]string{"id": conversationID})
		if err != nil {
			t.Fatal()
		}
	}))
	defer srv.Close()

	cli := NewClient(srv.Client(), srv.URL, fakeAPIKey, "model", "instructions")
	convID, err := cli.CreateConversation(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, convID, conversationID)
}
