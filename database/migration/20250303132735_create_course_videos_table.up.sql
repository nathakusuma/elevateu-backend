CREATE TABLE course_videos
(
    id          UUID PRIMARY KEY,
    course_id   UUID                     NOT NULL REFERENCES courses (id) ON DELETE CASCADE,
    title       VARCHAR(50)              NOT NULL,
    description VARCHAR(1000)            NOT NULL,
    duration    INT                      NOT NULL,
    is_free     BOOLEAN                  NOT NULL DEFAULT FALSE,
    "order"     INT                      NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX course_videos_order_index ON course_videos ("order");
