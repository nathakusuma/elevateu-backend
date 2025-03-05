CREATE TABLE course_video_progresses
(
    student_id    UUID    NOT NULL REFERENCES students (user_id) ON DELETE CASCADE,
    video_id      UUID    NOT NULL REFERENCES course_videos (id) ON DELETE CASCADE,
    last_position INT     NOT NULL DEFAULT 0,
    is_completed  BOOLEAN NOT NULL DEFAULT FALSE,
    PRIMARY KEY (student_id, video_id)
);
