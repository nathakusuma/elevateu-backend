CREATE TABLE mentoring_chats
(
    id         UUID PRIMARY KEY,
    student_id UUID                     NOT NULL REFERENCES students (user_id),
    mentor_id  UUID                     NOT NULL REFERENCES users (id),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_trial   BOOLEAN                  NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX mentoring_chats_student_id_mentor_id_key ON mentoring_chats (student_id, mentor_id);
CREATE INDEX mentoring_chats_student_id_idx ON mentoring_chats (student_id);
CREATE INDEX mentoring_chats_mentor_id_idx ON mentoring_chats (mentor_id);
