CREATE TABLE challenge_groups
(
    id              UUID PRIMARY KEY,
    category_id     UUID                     NOT NULL REFERENCES categories (id),
    title           VARCHAR(50)              NOT NULL,
    description     VARCHAR(1000)            NOT NULL,
    challenge_count INT                      NOT NULL DEFAULT 0,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX challenge_groups_title_idx ON challenge_groups USING gist (title gist_trgm_ops);
