package controller

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/vano2903/service-template/model"
	"github.com/vano2903/service-template/repo"
	"github.com/vano2903/service-template/repo/mock"
	"github.com/vano2903/service-template/services/logo"
)

var _ UserControllerer = new(UserController)

type UserController struct {
	repo repo.UserRepoer
	logo logo.LogoServicer
	l    logrus.Logger
}

func NewUserController(repo repo.UserRepoer, logo logo.LogoServicer, log logrus.Logger) *UserController {
	return &UserController{
		repo: repo,
		logo: logo,
		l:    log,
	}
}

func (c *UserController) CreateUser(firstName, lastName, email, password, role string) (int, error) {
	m := &model.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		Role:      role,
	}
	var err error
	m.Pfp, err = c.logo.GenerateLogo()
	if err != nil {
		c.l.Errorf("controller.CreateUser: unexpected error in logo.GenerateLogo: %v", err)
		return -1, errors.New("unexpected error when generating logo")
	}
	return c.repo.Create(m)
}

func (c *UserController) GetUser(id int) (*model.User, error) {
	return c.repo.Get(id)
}

func (c *UserController) GetAllUsers() ([]*model.User, error) {
	return c.repo.GetAll()
}

func (c *UserController) UpdateUser(requesterId int, u *model.User) error {
	requester, err := c.repo.Get(requesterId)
	if err != nil {
		//in this case we check if the error is a error not found and we log it
		//but not return it cause it could be something important that
		//the user should not know but us as developers should know
		re, ok := err.(*mock.ErrUserNotFound)
		if ok {
			//here we log the error and return a generic one
			c.l.Errorf("user with id %d not found", re.ID)
			return errors.New("user not found")
		} else {
			c.l.Errorf("controller.UpdateUser: unexpected error in repo.Get: %v", err)
			return errors.New("unexpected error")
		}
	}

	if requester.ID == u.ID || requester.Role == "admin" {
		err = c.repo.Update(u)
		if err != nil {
			re, ok := err.(*mock.ErrUserNotFound)
			if ok {
				c.l.Errorf("user with id %d not found", re.ID)
				return errors.New("user not found")

			} else if err == mock.ErrUserUnapdatable {
				return errors.New("user can't be updated")

			} else {
				c.l.Errorf("controller.UpdateUser: unexpected error in repo.Update: %v", err)
				return errors.New("unexpected error")
			}
		}
	}
	return nil
}

func (c *UserController) DeleteUser(requesterId, id int) error {
	requester, err := c.repo.Get(requesterId)
	if err != nil {
		re, ok := err.(*mock.ErrUserNotFound)
		if ok {
			c.l.Errorf("user with id %d not found", re.ID)
			return errors.New("user not found")
		} else {
			c.l.Errorf("controller.DeleteUser: unexpected error in repo.Get: %v", err)
			return errors.New("unexpected error")
		}
	}

	if requester.ID == id || requester.Role == "admin" {
		err = c.repo.Delete(id)
		if err != nil {
			re, ok := err.(*mock.ErrUserNotFound)
			if ok {
				c.l.Errorf("user with id %d not found", re.ID)
				return errors.New("user not found")
			} else {
				c.l.Errorf("controller.DeleteUser: unexpected error in repo.Delete: %v", err)
				return errors.New("unexpected error")
			}
		}
	}
	return nil
}

func (c *UserController) RegenerateLogo(id int) error {
	m, err := c.repo.Get(id)
	if err != nil {
		re, ok := err.(*mock.ErrUserNotFound)
		if ok {
			c.l.Errorf("user with id %d not found", re.ID)
			return errors.New("user not found")
		} else {
			c.l.Errorf("controller.RegenerateLogo: unexpected error in repo.Get: %v", err)
			return errors.New("unexpected error")
		}
	}

	m.Pfp, err = c.logo.GenerateLogo()
	if err != nil {
		c.l.Errorf("controller.RegenerateLogo: unexpected error in logo.GenerateLogo: %v", err)
		return errors.New("unexpected error")
	}

	err = c.repo.Update(m)
	if err != nil {
		re, ok := err.(*mock.ErrUserNotFound)
		if ok {
			c.l.Errorf("user with id %d not found", re.ID)
			return errors.New("user not found")
		} else if err == mock.ErrUserUnapdatable {
			return errors.New("user can't be updated")
		} else {
			c.l.Errorf("controller.RegenerateLogo: unexpected error in repo.Update: %v", err)
			return errors.New("unexpected error")
		}
	}
	return nil
}
