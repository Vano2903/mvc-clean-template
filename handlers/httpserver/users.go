package httpserver

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/vano2903/service-template/controller"
	"github.com/vano2903/service-template/repo/mock"
)

type (
	// Here we declare a user model that will be returned by the api
	// to any unauthorized user as some informations should be visible only to admins
	HttpUnauthenticatedUser struct {
		ID        int    `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Pfp       string `json:"pfp"`
		Email     string `json:"email"`
	}
	// We are not going to declare a model for the authorized request as we will just return the model

	userHttpHandler struct {
		e          *echo.Group
		controller *controller.User
		l          *logrus.Logger
	}
)

func NewUserHttpHandler(e *echo.Group, c *controller.User, l *logrus.Logger) *userHttpHandler {
	return &userHttpHandler{
		e:          e,
		controller: c,
		l:          l,
	}
}

// Registers only the routes and links functions
func (h *userHttpHandler) RegisterRoutes() {
	//user routes
	h.e.GET("/:id", h.GetUnauthorizedUser)
}

// @Summary		Get user from ID
// @Description	Get user from ID for unauthorized users
// @ID				getUnauthorizedUser
// @Tags			users
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	HttpSuccess{data=HttpUnauthenticatedUser,code=int,message=string}
// @Failure		500	{object}	HttpError
// @Failure		400	{object}	HttpError
// @Router			/user/{id} [get]
func (h *userHttpHandler) GetUnauthorizedUser(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		//we are not going to log this error as it is not an exception but a user error, we are just
		//telling the user what he did wrong but not saving the problem.
		//You could implement a request tracer that generates a request id and logs it, in that cause it would be more useful
		return respError(c, 400, "invalid id", fmt.Sprintf("id %q is not a valid id as it is not a number", idParam), "invalid_id")
	}

	user, err := h.controller.GetUser(id)
	if err != nil {
		_, ok := err.(*mock.ErrUserNotFound)
		if ok {
			return respError(c, 404, "user not found", fmt.Sprintf("user with id %d not found", id), "user_not_found")
		} else {
			return respError(c, 500, "unexpected error", fmt.Sprintf("unexpected error trying to retrive user %d", id), "unexpected_error")
		}
	}

	//converting the user to the http unauthorized user
	httpUser := HttpUnauthenticatedUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Pfp:       user.Pfp,
		Email:     user.Email,
	}

	return respSuccess(c, 200, "user succesfully retrived", httpUser)
}
