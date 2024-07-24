package main

type Exercise struct {
	ExerciseID          int      `json:"exerciseID"`
	ExerciseName        string   `json:"exerciseName"`
	MuscleNames         []string `json:"muscleNames"`
	MuscleGroups        []string `json:"muscleGroups"`
	ExerciseDescription string   `json:"exerciseDescription"`
	ReferenceVideo      string   `json:"referenceVideo"`
}

type CreateExerciseRequest struct {
	ExerciseName        string   `json:"exerciseName"`
	MuscleNames         []string `json:"muscleNames"`
	MuscleGroups        []string `json:"muscleGroups"`
	ExerciseDescription string   `json:"exerciseDescription"`
	ReferenceVideo      string   `json:"referenceVideo"`
}
