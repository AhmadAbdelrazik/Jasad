use jasad;

CREATE TABLE sessions(
    session_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    date DATE,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE TABLE workouts(
    workout_id  INT PRIMARY KEY AUTO_INCREMENT,
    session_id INT,
    exercise_id INT,
    reps INT,
    sets INT,
    weights FLOAT,
    FOREIGN KEY (exercise_id) REFERENCES exercises(exercise_id),
    FOREIGN KEY (session_id) REFERENCES exercises(exercise_id)
);

