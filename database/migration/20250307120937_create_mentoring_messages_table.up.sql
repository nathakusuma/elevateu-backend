CREATE TABLE mentoring_messages
(
    id         UUID PRIMARY KEY,
    chat_id    UUID                     NOT NULL REFERENCES mentoring_chats (id) ON DELETE CASCADE,
    sender_id  UUID                     NOT NULL REFERENCES users (id),
    message    VARCHAR(2000)            NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX mentoring_messages_chat_id_idx ON mentoring_messages (chat_id);
CREATE INDEX mentoring_messages_sender_id_idx ON mentoring_messages (sender_id);
