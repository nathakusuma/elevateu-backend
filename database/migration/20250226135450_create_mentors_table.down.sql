CREATE TABLE mentors
(
    user_id        UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    specialization VARCHAR(255)     NOT NULL,
    experience     VARCHAR(500)     NOT NULL,
    rating         FLOAT            NOT NULL DEFAULT 0,
    rating_count   INTEGER          NOT NULL DEFAULT 0,
    total_rating   DOUBLE PRECISION NOT NULL DEFAULT 0,
    price          INTEGER          NOT NULL,
    balance        INTEGER
);
