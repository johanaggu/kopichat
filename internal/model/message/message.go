package message

type Role string

const (
	BotRole Role = "bot"
	UserRol Role = "user"
)

// NewMessage creates a new Messages instance.
func NewMessage(id int, role Role, content, chatID string) *Message {
	return &Message{
		id:      id,
		role:    role,
		content: content,
		chatID:  chatID,
	}
}

// Messages represents a single message in a chat conversation.
type Message struct {
	id      int
	role    Role
	content string
	chatID  string
}

func (m *Message) ID() int {
	return m.id
}
func (m *Message) Role() Role {
	return m.role
}
func (m *Message) Content() string {
	return m.content
}

func (m *Message) ChatID() string {
	return m.chatID
}

func (m *Message) SetID(id int) {
	m.id = id
}
