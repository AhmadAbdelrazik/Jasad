CREATE TABLE IF NOT EXISTS users(
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	email VARCHAR(50) NOT NULL,
	role VARCHAR(30) NOT NULL,
	version INT NOT NULL DEFAULT 1
);
