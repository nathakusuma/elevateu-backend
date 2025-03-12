CREATE TABLE payments
(
    id         UUID PRIMARY KEY,
    user_id    UUID                     NOT NULL REFERENCES users (id),
    token      TEXT                     NOT NULL,
    amount     INT                      NOT NULL,
    title      VARCHAR(50)              NOT NULL,
    detail     VARCHAR(255),
    method     VARCHAR(50)              NOT NULL,
    status     VARCHAR(50)              NOT NULL,
    expired_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
