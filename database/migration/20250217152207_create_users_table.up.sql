CREATE TABLE users
(
    id            UUID PRIMARY KEY,
    name       VARCHAR(60) NOT NULL,
    email         VARCHAR(320) NOT NULL,
    password_hash CHAR(60) NOT NULL,
    role          VARCHAR(50)  NOT NULL DEFAULT 'student' CHECK ( role IN ('admin', 'mentor', 'student') ),
    avatar_url    TEXT,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX users_email_key ON users (email);
