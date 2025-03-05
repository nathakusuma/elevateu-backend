CREATE TABLE course_feedbacks
(
    id         UUID PRIMARY KEY,
    course_id  UUID                     NOT NULL REFERENCES courses (id) ON DELETE CASCADE,
    student_id UUID                     NOT NULL REFERENCES students (user_id) ON DELETE CASCADE,
    rating     DOUBLE PRECISION         NOT NULL,
    comment    VARCHAR(500)             NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX course_feedbacks_course_student_key ON course_feedbacks (course_id, student_id);
CREATE INDEX course_feedbacks_course_id_idx ON course_feedbacks (course_id);
