package model

import "errors"

type Muscle string

const (
	Shoulders  Muscle = "shoulder"
	Back              = "back"
	Traps             = "traps"
	Triceps           = "triceps"
	Biceps            = "biceps"
	Hands             = "hands"
	Lats              = "lats"
	LowerBack         = "lower back"
	Glutes            = "glutes"
	Hamstrings        = "hamstrings"
	Calves            = "calves"
	Quads             = "quads"
	Abdominals        = "abdominals"
	Obliques          = "obliques"
	Chest             = "chest"
)

func GetMuscle(s string) (Muscle, error) {
	switch s {
	case "shoulder":
		return Shoulders, nil
	case "back":
		return Back, nil
	case "traps":
		return Traps, nil
	case "triceps":
		return Triceps, nil
	case "biceps":
		return Biceps, nil
	case "hands":
		return Hands, nil
	case "lats":
		return Lats, nil
	case "lower back":
		return LowerBack, nil
	case "glutes":
		return Glutes, nil
	case "hamstrings":
		return Hamstrings, nil
	case "calves":
		return Calves, nil
	case "quads":
		return Quads, nil
	case "abdominals":
		return Abdominals, nil
	case "obliques":
		return Obliques, nil
	case "chest":
		return Chest, nil
	default:
		return "", errors.New("invalid muscle name")
	}
}
