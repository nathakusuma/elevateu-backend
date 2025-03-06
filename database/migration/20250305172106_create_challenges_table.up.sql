CREATE TABLE challenges
(
    id               UUID PRIMARY KEY,
    group_id         UUID                     NOT NULL REFERENCES challenge_groups (id) ON DELETE CASCADE,
    title            VARCHAR(50)              NOT NULL,
    subtitle         VARCHAR(100)             NOT NULL,
    description      VARCHAR(5000)            NOT NULL,
    difficulty       VARCHAR(20)              NOT NULL CHECK ( difficulty IN ('beginner', 'intermediate', 'advanced') ),
    is_free          BOOLEAN                  NOT NULL DEFAULT FALSE,
    submission_count BIGINT                   NOT NULL DEFAULT 0,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX challenges_group_id_idx ON challenges (group_id);
CREATE INDEX challenges_difficulty_idx ON challenges (difficulty);
