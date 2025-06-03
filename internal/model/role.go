package model

import "errors"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser       = "user"
)

func GetRole(role string) (Role, error) {
	switch role {
	case "admin":
		return RoleAdmin, nil
	case "user":
		return RoleUser, nil
	default:
		return Role(""), errors.New("invalid role")
	}
}
