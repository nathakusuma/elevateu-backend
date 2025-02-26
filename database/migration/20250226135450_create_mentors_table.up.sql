CREATE TABLE mentors
(
    user_id        UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    specialization VARCHAR(255)     NOT NULL,
    experience     VARCHAR(1000)    NOT NULL,
    rating         DOUBLE PRECISION NOT NULL DEFAULT 0,
    rating_count   INTEGER          NOT NULL DEFAULT 0,
    rating_total   DOUBLE PRECISION NOT NULL DEFAULT 0,
    price          INTEGER          NOT NULL,
    balance        INTEGER          NOT NULL DEFAULT 0
);
