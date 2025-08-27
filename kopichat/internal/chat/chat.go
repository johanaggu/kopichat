package chat

// NewChat creates a new Chat entity with the given id and externalID.
func NewChat(id, externalID string) *Chat {
	return &Chat{
		id:     id,
		externalID: externalID,
	}
}

// Chat represents a chat entity, this object contains a unique ID and the chat ID from the external API.
// in this case the chat ID is from the OpenAI API.
type Chat struct {
	// ID is the unique identifier for the chat entity.
	id string
	// externalID is the identifier used by the Open AI to reference the chat.
	externalID string
}

// ID returns the unique identifier of the chat entity.
func (c *Chat) ID() string {
	return c.id
}

// ChatID returns the identifier used by the Open AI to reference the chat.
func (c *Chat) ExternalID() string {
	return c.externalID
}
