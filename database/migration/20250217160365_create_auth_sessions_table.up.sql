CREATE TABLE auth_sessions
(
    token      CHAR(32) PRIMARY KEY,
    user_id    UUID REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX auth_sessions_expires_at_idx ON auth_sessions (expires_at);
CREATE UNIQUE INDEX auth_sessions_user_id_key ON auth_sessions (user_id);
