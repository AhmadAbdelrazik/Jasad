package main

import (
	"fmt"

	"github.com/AhmadAbdelrazik/jasad/internal/validator"
)


func CheckExerciseForm(form *ExerciseForm) {
	
	form.CheckField(validator.NotBlank(form.ExerciseName), "exercise_name", validator.MessageEmptyField)
	form.CheckField(validator.Chars(form.ExerciseName, 5, 50), "exercise_name", validator.MessageFieldLength(5, 50))

	form.CheckField(validator.NotBlank(form.ExerciseDescription), "exercise_description", validator.MessageEmptyField)
	form.CheckField(validator.Chars(form.ExerciseDescription ,5, 500), "exercise_description", validator.MessageFieldLength(5, 500))

	form.CheckField(len(form.Muscles) > 0,"muscles", validator.MessageEmptyArray)

	for _, muscle := range form.Muscles {
		form.CheckField(validator.NotBlank(muscle),fmt.Sprintf("muscle_%v", muscle), validator.MessageEmptyField)
		form.CheckField(validator.MaxChars(muscle, 20),fmt.Sprintf("muscle_%v", muscle) , validator.MessageFieldLength(2, 20))

	}
	
	form.CheckField(validator.NotBlank(form.ReferenceVideo), "reference_video", validator.MessageEmptyField)
	form.CheckField(validator.Chars(form.ReferenceVideo, 5, 50), "reference_video", validator.MessageFieldLength(5, 50))
}