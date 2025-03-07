CREATE TABLE challenge_submission_feedbacks
(
    submission_id UUID PRIMARY KEY REFERENCES challenge_submissions (id),
    mentor_id     UUID                     NOT NULL REFERENCES users (id),
    score         INT                      NOT NULL CHECK ( score >= 0 AND score <= 100 ),
    feedback      VARCHAR(1000)            NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX challenge_submission_feedbacks_mentor_id_idx ON challenge_submission_feedbacks (mentor_id);
