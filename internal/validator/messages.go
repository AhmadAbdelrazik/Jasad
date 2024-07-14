package validator

import "fmt"

var (
	MessageEmptyField = "this field can't be empty"
	MessageEmptyArray = "this list should contain at least one item"
)

func MessageFieldLength(min, max int) string {
	return fmt.Sprintf("this field must be between %v and %v characters", min, max)
}
