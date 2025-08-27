CREATE TABLE chats (
  id    INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  api_id VARCHAR(255) NOT NULL,
);

CREATE INDEX idx_chats_api_id ON chats (api_id);