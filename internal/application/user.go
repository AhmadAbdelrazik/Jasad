package application

import (
	"errors"
	"net/http"

	"github.com/ahmadabdelrazik/jasad/internal/model"
)

func (app *Application) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		BadRequestResponse(w, r, err)
		return
	}

	authUser, ok := getUser(r)
	if !ok { // if there is no user in context
		UnauthorizedResponse(w, r)
		return
		// if the user is not asking for his info and isn't an admin.
	} else if authUser.ID != int(id) && authUser.Role != model.RoleAdmin {
		UnauthorizedResponse(w, r)
		return
	}

	user, err := app.models.Users.GetByID(int(id))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			NotFoundResponse(w, r)
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	output := struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	}{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": output}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}

func (app *Application) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.Users.GetAll()
	if err != nil {
		ServerErrorResponse(w, r, err)
		return
	}

	//NOTE: Can use a special output struct to protect any sensetive info
	//later in the project

	err = app.writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}
