CREATE TABLE course_enrollments
(
    course_id         UUID                     NOT NULL REFERENCES courses (id) ON DELETE CASCADE,
    student_id        UUID                     NOT NULL REFERENCES students (user_id) ON DELETE CASCADE,
    content_completed INT                      NOT NULL DEFAULT 0,
    is_completed      BOOLEAN                  NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_accessed_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (course_id, student_id)
);

CREATE INDEX course_enrollments_student_lastaccessed_course_idx
    ON course_enrollments (student_id, last_accessed_at DESC, course_id);
