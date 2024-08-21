use jasad;

ALTER TABLE workouts
DROP FOREIGN KEY workouts_ibfk_2;

ALTER TABLE workouts
ADD CONSTRAINT workouts_ibfk_2
FOREIGN KEY (session_id) REFERENCES sessions(session_id);
