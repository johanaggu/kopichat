CREATE TABLE messages (
  `id`              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  `chat_id`         CHAR(36) NOT NULL,
  `content`         TEXT NOT NULL,
  `role`            VARCHAR(20) NOT NULL,
  `created_at`      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_chat_messages FOREIGN KEY (`chat_id`)
  REFERENCES chats(id)
  ON DELETE CASCADE
);

CREATE INDEX idx_chat_messages ON messages(chat_id);