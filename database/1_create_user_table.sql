use jasad;

CREATE TABLE users (
    user_id INT PRIMARY KEY AUTO_INCREMENT,
    user_name VARCHAR(32) UNIQUE,
    password VARCHAR(90),
    role VARCHAR(20),
    created_at DATETIME
);

CREATE TABLE exercises (
    exercise_id INT PRIMARY KEY AUTO_INCREMENT,
    exercise_name VARCHAR(30) UNIQUE,
    exercise_description TEXT,
    reference_video VARCHAR(100)
);

CREATE TABLE muscles_exercises (
    exercise_id INT,
    muscle_name VARCHAR(30),
    muscle_group VARCHAR(30),
    FOREIGN KEY (exercise_id) REFERENCES exercises(exercise_id)
);
