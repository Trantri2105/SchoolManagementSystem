CREATE TABLE IF NOT EXISTS users(
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    date_of_birth DATE,
    gender TEXT,
    email TEXT UNIQUE,
    identity_number TEXT,
    phone_number TEXT,
    address TEXT,
    password TEXT,
    role TEXT
);

CREATE TABLE IF NOT EXISTS students(
    id TEXT PRIMARY KEY REFERENCES users(id),
    school_year TEXT,
    major TEXT
);

CREATE TABLE IF NOT EXISTS teachers(
    id TEXT PRIMARY KEY REFERENCES users(id),
    academic_qualification TEXT NOT NULL,
    department TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS subjects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    number_of_credit INT NOT NULL,
    major TEXT
);

CREATE TABLE IF NOT EXISTS courses (
    id TEXT PRIMARY KEY,
    teacher_id TEXT REFERENCES teachers(id) ON DELETE CASCADE,
    subject_id TEXT REFERENCES subjects(id) ON DELETE CASCADE,
    semester_number INT NOT NULL,
    academic_year TEXT NOT NULL,
    status TEXT,
    capacity INT,
    size INT
);

CREATE TABLE IF NOT EXISTS course_schedules (
    id SERIAL PRIMARY KEY,
    course_id TEXT REFERENCES courses(id) ON DELETE CASCADE,
    room TEXT,
    start_time INT,
    end_time INT
);

CREATE TABLE IF NOT EXISTS course_registrations (
    id SERIAL PRIMARY KEY,
    course_id TEXT REFERENCES courses(id) ON DELETE CASCADE,
    student_id TEXT REFERENCES students(id) ON DELETE CASCADE,
    UNIQUE (course_id,student_id)
);

CREATE TABLE IF NOT EXISTS component_scores(
    id SERIAL PRIMARY KEY,
    course_id TEXT REFERENCES courses(id) ON DELETE CASCADE,
    name TEXT,
    score_weight DOUBLE PRECISION,
    score DOUBLE PRECISION
);

INSERT INTO users (
    id, name, date_of_birth, gender, email, identity_number, phone_number, address, password, role
) VALUES (
             'admin001',
             'Admin User',
             '1980-01-01',
             'Male',
             'admin@example.com',
             '123456789',
             '1234567890',
             '123 Admin Street, Admin City',
             '$2a$04$PFxZjH2Jwb2aL0yslR8N.eo8Cw.PEcdFjB68a7d3tBmImq3yAtR1S',
             'Admin'
         );






