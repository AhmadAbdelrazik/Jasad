package main


type Muscle struct {
	MuscleName  string `json:"muscleName" validate:"required,max=20"`
	MuscleGroup string `json:"muscleGroup" validate:"required,max=20"`
}

type Exercise struct {
	ExerciseID          int      `json:"exerciseID"`
	ExerciseName        string   `json:"exerciseName"`
	ExerciseDescription string   `json:"exerciseDescription"`
	ReferenceVideo      string   `json:"referenceVideo"`
	Muscles             []Muscle `json:"muscles"`
}

type CreateExerciseRequest struct {
	ExerciseName        string   `json:"exerciseName" validate:"required,min=3,max=30"`
	ExerciseDescription string   `json:"exerciseDescription" validate:"required,max=300"`
	ReferenceVideo      string   `json:"referenceVideo" validate:"required,url"`
	Muscles             []Muscle `json:"muscles"`
}
