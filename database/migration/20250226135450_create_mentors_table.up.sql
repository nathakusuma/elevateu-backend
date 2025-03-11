CREATE TABLE mentors
(
    user_id        UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    address        VARCHAR(255) NOT NULL,
    specialization VARCHAR(50)  NOT NULL,
    current_job    VARCHAR(50)  NOT NULL,
    company        VARCHAR(50)  NOT NULL,
    bio            VARCHAR(1000),
    gender         VARCHAR(10)  NOT NULL CHECK (gender IN ('male', 'female')),
    rating         DOUBLE PRECISION NOT NULL DEFAULT 0,
    rating_count   INTEGER          NOT NULL DEFAULT 0,
    rating_total   DOUBLE PRECISION NOT NULL DEFAULT 0,
    price          INTEGER      NOT NULL DEFAULT 40000,
    balance        INTEGER          NOT NULL DEFAULT 0
);
