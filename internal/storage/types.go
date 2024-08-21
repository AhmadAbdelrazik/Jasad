package storage

import "time"

type Muscle struct {
	MuscleName  string `json:"muscleName" validate:"required,max=20"`
	MuscleGroup string `json:"muscleGroup" validate:"required,max=20"`
}

type Exercise struct {
	ExerciseID          int      `json:"exerciseID" validate:"required,gt=0"`
	ExerciseName        string   `json:"exerciseName" validate:"required,min=3,max=30"`
	ExerciseDescription string   `json:"exerciseDescription" validate:"required,max=300"`
	ReferenceVideo      string   `json:"referenceVideo" validate:"required,url"`
	Muscles             []Muscle `json:"muscles"`
}

type ExerciseUpdateRequest struct {
	ExerciseID          int      `json:"exerciseID" validate:"required,gt=0"`
	ExerciseName        string   `json:"exerciseName" validate:"required,min=3,max=30"`
	ExerciseDescription string   `json:"exerciseDescription" validate:"required,max=300"`
	ReferenceVideo      string   `json:"referenceVideo" validate:"required,url"`
	Muscles             []Muscle `json:"muscles"`
}

type ExerciseCreateRequest struct {
	ExerciseName        string   `json:"exerciseName" validate:"required,min=3,max=30"`
	ExerciseDescription string   `json:"exerciseDescription" validate:"required,max=300"`
	ReferenceVideo      string   `json:"referenceVideo" validate:"required,url"`
	Muscles             []Muscle `json:"muscles"`
}

// This struct shall not be used by any endpoint as a response
type User struct {
	UserName  string    `json:"username"`
	UserID    int       `json:"userID"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserRequest struct {
	UserName string `json:"username" validate:"required,min=8,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type UserJWT struct {
	UserName string `json:"username"`
	Role     string `json:"role"`
}

type Workout struct {
	WorkoutID  int     `json:"workoutID" validate:"required"`
	ExerciseID int     `json:"exerciseID" validate:"required"`
	Reps       int     `json:"reps" validate:"required"`
	Sets       int     `json:"sets" validate:"required"`
	Weights    float32 `json:"weights" validate:"required"`
}

type Session struct {
	SessionID int       `json:"sessionID"`
	Workouts  []Workout `json:"workouts"`
	Date      time.Time `json:"date"`
	UserID    int       `json:"userID"`
}

type WorkoutCreateRequest struct {
	Workouts []Workout `json:"workouts" validate:"required"`
	Date     time.Time `json:"date" validate:"required"`
}

type SessionResponse struct {
	SessionID int       `json:"sessionID"`
	Date      time.Time `json:"date"`
}

type WorkoutGetResponse struct {
	Username string            `json:"username"`
	Sessions []SessionResponse `json:"sessions"`
}
