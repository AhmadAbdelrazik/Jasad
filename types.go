package main

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

type UpdateExerciseRequest struct {
	ExerciseID          int      `json:"exerciseID" validate:"required,gt=0"`
	ExerciseName        string   `json:"exerciseName" validate:"required,min=3,max=30"`
	ExerciseDescription string   `json:"exerciseDescription" validate:"required,max=300"`
	ReferenceVideo      string   `json:"referenceVideo" validate:"required,url"`
	Muscles             []Muscle `json:"muscles"`
}

type CreateExerciseRequest struct {
	ExerciseName        string   `json:"exerciseName" validate:"required,min=3,max=30"`
	ExerciseDescription string   `json:"exerciseDescription" validate:"required,max=300"`
	ReferenceVideo      string   `json:"referenceVideo" validate:"required,url"`
	Muscles             []Muscle `json:"muscles"`
}

// This struct shall not be used by any endpoint as a response
type User struct {
	UserName  string    `json:"username"`
	UserID    string    `json:"userID"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateUserRequest struct {
	UserName string `json:"username" validate:"required,min=8,max=32"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}
