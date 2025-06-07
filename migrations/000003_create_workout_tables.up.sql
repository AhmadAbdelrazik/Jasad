CREATE TABLE IF NOT EXISTS workouts(
	id SERIAL PRIMARY KEY,
	owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	name VARCHAR(50) NOT NULL,
	version INT NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS workouts_exercises(
	id SERIAL PRIMARY KEY,
	exercise_id INT NOT NULL REFERENCES exercises(id),
	workout_id INT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
	exercise_order INT NOT NULL,
	sets INT NOT NULL,
	reps INT NOT NULL,
	weights REAL NOT NULL,
	rest_after INT NOT NULL,
	done bool NOT NULL,
	version INT NOT NULL DEFAULT 1
);
