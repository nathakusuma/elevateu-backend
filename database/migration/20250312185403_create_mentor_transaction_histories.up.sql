CREATE TABLE mentor_transaction_histories
(
    id         UUID PRIMARY KEY,
    mentor_id  UUID                     NOT NULL REFERENCES users (id),
    title      VARCHAR(50)              NOT NULL,
    detail     VARCHAR(255),
    amount     INT                      NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
