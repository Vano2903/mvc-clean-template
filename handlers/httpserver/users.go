package httpserver

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/vano2903/service-template/controller"
	"github.com/vano2903/service-template/model"
	"github.com/vano2903/service-template/pkg/jwt"
	"github.com/vano2903/service-template/repo/mock"
)

const (
	bearerHeaderLength = 7
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

	HttpNewUserPost struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	HttpUpdateUserPost struct {
		ID        int    `json:"id,omitempty" validate:"optional"`
		FirstName string `json:"first_name,omitempty" validate:"optional"`
		LastName  string `json:"last_name,omitempty" validate:"optional"`
		Email     string `json:"email,omitempty" validate:"optional"`
		Password  string `json:"password,omitempty" validate:"optional"`
	}

	HttpLoginUserPost struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	userHttpHandler struct {
		e          *echo.Group
		controller *controller.User
		l          *logrus.Logger
		j          *jwt.JWThandler
	}
)

func NewUserHttpHandler(e *echo.Group, c *controller.User, l *logrus.Logger, jwtHandler *jwt.JWThandler) *userHttpHandler {
	return &userHttpHandler{
		e:          e,
		controller: c,
		l:          l,
		j:          jwtHandler,
	}
}

// Registers only the routes and links functions
func (h *userHttpHandler) RegisterRoutes() {
	//user routes
	h.e.GET("/:id", h.GetUnauthorizedUser)
	h.e.GET("/all", h.GetAllUnauthorizedUsers)
	h.e.POST("/register", h.CreateNewUser)
	h.e.POST("/login", h.LoginUser)

	h.e.GET("/me", h.GetUserInfo, h.jwtHeaderCheckerMiddleware)
	h.e.GET("/update", h.UpdateUser, h.jwtHeaderCheckerMiddleware)
}

// @Summary		Get user from ID
// @Description	Get user from ID for unauthorized users
// @ID				getUnauthorizedUser
// @Tags			users
// @Produce		json
// @Param			id	path		int	true	"User ID"
// @Success		200	{object}	HttpSuccess{data=HttpUnauthenticatedUser,code=int,message=string}
// @Failure		400	{object}	HttpError
// @Failure		404	{object}	HttpError
// @Failure		500	{object}	HttpError
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

// @Summary		Get all user
// @Description	Get all user for an unauthorized user
// @ID				getAllUnauthorizedUser
// @Tags			users
// @Produce		json
// @Success		200	{object}	HttpSuccess{data=[]HttpUnauthenticatedUser,code=int,message=string}
// @Failure		404	{object}	HttpError
// @Failure		500	{object}	HttpError
// @Router			/user/all [get]
func (h *userHttpHandler) GetAllUnauthorizedUsers(c echo.Context) error {
	users := h.controller.GetAllUsers()
	if len(users) == 0 {
		return respError(c, 404, "no users found", "no users were found for this unauthorized access", "no_users_found")
	}
	unauthUser := make([]HttpUnauthenticatedUser, len(users))
	for _, u := range users {
		unauthUser = append(unauthUser, HttpUnauthenticatedUser{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Pfp:       u.Pfp,
			Email:     u.Email,
		})
	}

	return respSuccess(c, 200, "all users succesfully retrived", unauthUser)
}

// @Summary		Register a new user
// @Description	Register a new user
// @ID				CreateNewUser
// @Tags			users
// @Produce		json
// @Param			account	body		HttpNewUserPost	true	"User Informations"
// @Success		200		{object}	HttpSuccess{data=httpserver.CreateNewUser.HttpNewUserPostResponse,code=int,message=string}
// @Failure		400		{object}	HttpError
// @Failure		500		{object}	HttpError
// @Router			/user/register [POST]
func (h *userHttpHandler) CreateNewUser(c echo.Context) error {
	body := HttpNewUserPost{}
	if err := c.Bind(&body); err != nil {
		return respError(c, 400, "invalid body", fmt.Sprintf("invalid body: %v", err), "invalid_body")
	}

	newUserID, err := h.controller.CreateUser(body.FirstName, body.LastName, body.Email, body.Password, model.RoleUser)
	if err != nil {
		if err == controller.ErrUserAlreadyExists {
			return respError(c, 400, "user already exists", fmt.Sprintf("user with email %s already exists", body.Email), "user_already_exists")
		} else {
			return respError(c, 500, "unexpected error", fmt.Sprintf("unexpected error trying to create user %s", body.Email), "unexpected_error")
		}
	}

	type HttpNewUserPostResponse struct {
		ID int `json:"id"`
	}

	return respSuccess(c, 200, "user succesfully created", HttpNewUserPostResponse{ID: newUserID})
}

// @Summary		Login
// @Description	Login user given email and password
// @ID				LoginUser
// @Tags			users
// @Produce		json
// @Param			credentials	body		HttpLoginUserPost	true	"email and password"
// @Success		200			{object}	HttpSuccess{data=httpserver.LoginUser.HttpLoginUserPostResponse,code=int,message=string}
// @Failure		400			{object}	HttpError
// @Failure		401			{object}	HttpError
// @Failure		404			{object}	HttpError
// @Failure		500			{object}	HttpError
// @Router			/user/login [POST]
func (h *userHttpHandler) LoginUser(c echo.Context) error {
	body := HttpLoginUserPost{}
	if err := c.Bind(&body); err != nil {
		return respError(c, 400, "invalid body", fmt.Sprintf("invalid body: %v", err), "invalid_body")
	}

	id, err := h.controller.CheckCredentials(body.Email, body.Password)
	if err != nil {
		if err == controller.ErrUserNotFound {
			return respError(c, 404, "user not found", fmt.Sprintf("there is no user with %s as email", body.Email), "user_not_found")
		} else if err == controller.ErrWrongPassword {
			return respError(c, 401, "wrong password", "the password is not valid, check if spelled right", "wrong_password")
		} else {
			return respError(c, 500, "unexpected error", fmt.Sprintf("unexpected error trying to login user %s", body.Email), "unexpected_error")
		}
	}

	user, _ := h.controller.GetUser(id)

	jwtString, err := h.j.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		h.l.Errorf("unexpected error trying to sign jwt token for user %s: %v", body.Email, err)
		return respError(c, 500, "unexpected error", "unexpected error trying to generate your login token", "unexpected_error")
	}
	h.l.Debugf("token generated for user %s: %s", body.Email, jwtString)

	type HttpLoginUserPostResponse struct {
		Token string `json:"token"`
	}

	return respSuccess(c, 200, "user succesfully logged in", HttpLoginUserPostResponse{Token: jwtString})
}

// @Summary		Get user info
// @Description	Get authenticated user info from jwt
// @ID				GetUserInfo
// @Tags			users
// @Produce		json
// @Param Authorization header string  true "jwt token"     default(Bearer xxx.xxx.xxx)
// @Success		200			{object}	HttpSuccess{data=model.User,code=int,message=string}
// @Failure		401			{object}	HttpError
// @Failure		404			{object}	HttpError
// @Failure		500			{object}	HttpError
// @Router			/user/me [GET]
func (h *userHttpHandler) GetUserInfo(c echo.Context) error {
	//it wont panic because the middleware already checked it
	authHeader := c.Request().Header.Get("Authorization")[bearerHeaderLength:]

	//we do not do the full check as the middleware already did it
	//we just get the claims and handle the error
	claims, err := h.j.ValidateToken(authHeader)
	if err != nil {
		return respError(c, 401, "invalid token", "invalid token", "invalid_token")
	}

	user, err := h.controller.GetUser(claims.UserId)
	if err != nil {
		if err == controller.ErrUserNotFound {
			return respError(c, 404, "user not found", fmt.Sprintf("there is no user with %d as id", claims.UserId), "user_not_found")
		} else {
			return respError(c, 500, "unexpected error", fmt.Sprintf("unexpected error trying to retrive user %d", claims.UserId), "unexpected_error")
		}
	}

	return respSuccess(c, 200, "user succesfully retrived", user)
}

// @Summary		Update user
// @Description	Update user info from jwt, if you are an admin you can update any user given the id
// @Description A normal user can only update himself, an admin can update any user
// @Description You don't have to send all the fields, only the ones you want to update with the new values
// @Description If the ID is not specified the user the update will be applied to the requesting user (only for admins, normal users can't update other users)
// @ID				UpdateUser
// @Tags			users
// @Produce		json
// @Param Authorization header string  true "jwt token"     default(Bearer xxx.xxx.xxx)
// @Param		user_info	body		HttpUpdateUserPost	true	"users information to update"
// @Success		200			{object}	HttpSuccess{code=int,message=string}
// @Failure		400			{object}	HttpError
// @Failure		401			{object}	HttpError
// @Failure		404			{object}	HttpError
// @Failure		500			{object}	HttpError
// @Router			/user/update [POST]
func (h *userHttpHandler) UpdateUser(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")[bearerHeaderLength:]
	claims, err := h.j.ValidateToken(authHeader)
	if err != nil {
		return respError(c, 401, "invalid token", "invalid token", "invalid_token")
	}

	body := HttpUpdateUserPost{}
	if err := c.Bind(&body); err != nil {
		return respError(c, 400, "invalid body", fmt.Sprintf("invalid body: %v", err), "invalid_body")
	}

	toUpdateID := claims.UserId
	if body.ID != 0 {
		toUpdateID = body.ID
	}

	toUpdate, err := h.controller.GetUser(toUpdateID)
	if err != nil {
		if err == controller.ErrUserNotFound {
			_, err := h.controller.GetUser(claims.UserId)
			if err != nil {
				//this is an extreme case, if the user deleted the account the related jwt should be deleted aswell by putting it in a blacklist
				return respError(c, 404, "user not found", "the jwt references an user that does not exist, maybe the user deleted the account", "user_not_found")
			} else {
				return respError(c, 404, "user not found", fmt.Sprintf("there is no user with %d as id", toUpdateID), "user_not_found")
			}
		} else {
			return respError(c, 500, "unexpected error", fmt.Sprintf("unexpected error trying to retrive user %d", toUpdateID), "unexpected_error")
		}
	}

	if body.FirstName != "" {
		toUpdate.FirstName = body.FirstName
	} else if body.LastName != "" {
		toUpdate.LastName = body.LastName
	} else if body.Email != "" {
		toUpdate.Email = body.Email
	} else if body.Password != "" {
		toUpdate.Password = body.Password
	}

	if err := h.controller.UpdateUser(claims.UserId, toUpdate); err != nil {
		if err == controller.ErrUnupdatableUser {
			return respError(c, 400, "unupdatable user", "the user you are trying to update is not updatable", "unupdatable_user")
		} else {
			//we are excluding user not found because we already checked it
			return respError(c, 500, "unexpected error", fmt.Sprintf("unexpected error trying to update user %d", toUpdateID), "unexpected_error")
		}
	}
	return respSuccess(c, 200, "user succesfully updated")
}

// @Summary Regenerate user's pfp
// @Description Automatically regnerate user's profile picture
// @Description Normal user do not need to specify anything, admins can specify a userid to update
// @ID				RegeneratePfpUrl
// @Tags			users
// @Produce		json
// @Param Authorization header string  true "jwt token"     default(Bearer xxx.xxx.xxx)
// @Param		userid	path		int	false	"id of the user to update"
// @Success		200			{object}	HttpSuccess{data=httpserver.RegeneratePfpUrl.HttpNewPfp,code=int,message=string}
// @Failure		400			{object}	HttpError
// @Failure		401			{object}	HttpError
// @Failure		404			{object}	HttpError
// @Failure		500			{object}	HttpError
// @Router			/user/update [POST]
func (h *userHttpHandler) RegeneratePfpUrl(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")[bearerHeaderLength:]
	claims, err := h.j.ValidateToken(authHeader)
	if err != nil {
		return respError(c, 401, "invalid token", "invalid token", "invalid_token")
	}
	u, err := h.controller.GetUser(claims.UserId)
	if err != nil {
		return respError(c, 404, "user not found", "the jwt references an user that does not exist, maybe the user deleted the account", "user_not_found")
	}
	userIdToUpdate := claims.UserId

	if c.Param("userid") != "" {
		if u.Role != model.RoleAdmin {
			return respError(c, 401, "unauthorized", "only admins can update other users", "unauthorized_update")
		} else {
			userIdToUpdate, err = strconv.Atoi(c.Param("userid"))
			if err != nil {
				return respError(c, 400, "invalid userid", "the userid must be an integer", "invalid_userid")
			}
		}
	}

	if err := h.controller.RegeneratePfp(userIdToUpdate); err != nil {
		if err == controller.ErrUserNotFound {
			return respError(c, 404, "user not found", fmt.Sprintf("there is no user with %d as id", userIdToUpdate), "user_not_found")
		} else if err == controller.ErrUnupdatableUser {
			return respError(c, 400, "unupdatable user", "the user you are trying to update is not updatable", "unupdatable_user")
		} else {
			return respError(c, 500, "unexpected error", fmt.Sprintf("unexpected error trying to update user %d", userIdToUpdate), "unexpected_error")
		}
	}
	type HttpNewPfp struct {
		NewPfp string `json:"new_pfp"`
	}
	u, _ = h.controller.GetUser(userIdToUpdate)
	return respSuccess(c, 200, "user succesfully updated with a new profile picture", HttpNewPfp{NewPfp: u.Pfp})
}
