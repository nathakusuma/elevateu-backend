CREATE TABLE course_material_progresses
(
    student_id  UUID NOT NULL REFERENCES students (user_id) ON DELETE CASCADE,
    material_id UUID NOT NULL REFERENCES course_materials (id) ON DELETE CASCADE,
    PRIMARY KEY (student_id, material_id)
);
