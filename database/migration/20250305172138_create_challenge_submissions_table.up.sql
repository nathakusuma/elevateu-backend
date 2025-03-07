CREATE TABLE challenge_submissions
(
    id           UUID PRIMARY KEY,
    challenge_id UUID                     NOT NULL REFERENCES challenges (id) ON DELETE CASCADE,
    student_id   UUID                     NOT NULL REFERENCES students (user_id) ON DELETE CASCADE,
    url          VARCHAR(500)             NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX challenge_submissions_challenge_id_student_id_key ON challenge_submissions (challenge_id, student_id);
