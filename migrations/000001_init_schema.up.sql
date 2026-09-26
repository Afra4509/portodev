-- Create students table
CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    nrp VARCHAR(20) NOT NULL UNIQUE,
    nama VARCHAR(100) NOT NULL,
    program_studi VARCHAR(100) NOT NULL,
    fakultas VARCHAR(100) NOT NULL,
    angkatan INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index on student NRP for fast lookup
CREATE INDEX IF NOT EXISTS idx_students_nrp ON students(nrp);

-- Create courses table
CREATE TABLE IF NOT EXISTS courses (
    id SERIAL PRIMARY KEY,
    kode_mata_kuliah VARCHAR(20) NOT NULL UNIQUE,
    nama_mata_kuliah VARCHAR(100) NOT NULL,
    sks INT NOT NULL CHECK (sks >= 1 AND sks <= 6),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Index on course code
CREATE INDEX IF NOT EXISTS idx_courses_kode ON courses(kode_mata_kuliah);

-- Create grades table
CREATE TABLE IF NOT EXISTS grades (
    id SERIAL PRIMARY KEY,
    student_id INT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    course_id INT NOT NULL REFERENCES courses(id) ON DELETE RESTRICT,
    semester INT NOT NULL CHECK (semester >= 1 AND semester <= 14),
    nilai_huruf VARCHAR(5) NOT NULL CHECK (nilai_huruf IN ('A', 'AB', 'B', 'BC', 'C', 'D', 'E')),
    nilai_angka NUMERIC(4,2) NOT NULL CHECK (nilai_angka >= 0.00 AND nilai_angka <= 4.00),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_student_course_semester UNIQUE (student_id, course_id, semester)
);

-- Indexes on foreign keys and compound queries
CREATE INDEX IF NOT EXISTS idx_grades_student_id ON grades(student_id);
CREATE INDEX IF NOT EXISTS idx_grades_course_id ON grades(course_id);
CREATE INDEX IF NOT EXISTS idx_grades_student_semester ON grades(student_id, semester);
