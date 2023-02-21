package controller

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/vano2903/service-template/model"
	"github.com/vano2903/service-template/providers/logo"
	"github.com/vano2903/service-template/repo"
	"github.com/vano2903/service-template/repo/mock"
)

var _ UserControllerer = new(User)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrWrongPassword     = errors.New("wrong password")
	ErrUnexpected        = errors.New("unexpected error")
	ErrUnupdatableUser   = errors.New("user can't be updated")
)

type User struct {
	repo repo.UserRepoer
	logo logo.LogoServicer
	l    *logrus.Logger
}

func NewUserController(repo repo.UserRepoer, logo logo.LogoServicer, log *logrus.Logger) *User {
	return &User{
		repo: repo,
		logo: logo,
		l:    log,
	}
}

func (c *User) CreateUser(firstName, lastName, email, password, role string) (int, error) {

	//here we check if the user already exists
	u, err := c.repo.GetByEmail(email)
	if err == nil {
		return u.ID, ErrUserAlreadyExists
	}

	m := &model.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		Role:      role,
	}

	m.Pfp, err = c.logo.GenerateLogo()
	if err != nil {
		c.l.Errorf("controller.CreateUser: unexpected error in logo.GenerateLogo: %v", err)
		return u.ID, errors.New("unexpected error when generating logo")
	}
	return c.repo.Create(m)
}

func (c *User) GetUser(id int) (*model.User, error) {
	return c.repo.Get(id)
}

func (c *User) GetAllUsers() []*model.User {
	return c.repo.GetAll()
}

func (c *User) UpdateUser(requesterId int, u *model.User) error {
	requester, err := c.repo.Get(requesterId)
	if err != nil {
		//in this case we check if the error is a error not found and we log it
		//but not return it cause it could be something important that
		//the user should not know but us as developers should know
		re, ok := err.(*mock.ErrUserNotFound)
		if ok {
			//here we log the error and return a generic one
			c.l.Errorf("update requester with id %d not found", re.ID)
			return ErrUserNotFound
		} else {
			c.l.Errorf("controller.UpdateUser: unexpected error in repo.Get: %v", err)
			return ErrUnexpected
		}
	}

	if u.ID <= 0 {
		return errors.New("missing id from user to update")
	}

	if requester.ID == u.ID || requester.Role == model.RoleAdmin {
		err = c.repo.Update(u)
		if err != nil {
			re, ok := err.(*mock.ErrUserNotFound)
			if ok {
				c.l.Errorf("user to update with id %d not found", re.ID)
				return ErrUserNotFound
			} else if err == mock.ErrUserUnapdatable {
				return ErrUnupdatableUser
			} else {
				c.l.Errorf("controller.UpdateUser: unexpected error in repo.Update: %v", err)
				return ErrUnexpected
			}
		}
	}
	return nil
}

func (c *User) DeleteUser(requesterId, id int) error {
	requester, err := c.repo.Get(requesterId)
	if err != nil {
		re, ok := err.(*mock.ErrUserNotFound)
		if ok {
			c.l.Errorf("user with id %d not found", re.ID)
			return ErrUserNotFound
		} else {
			c.l.Errorf("controller.DeleteUser: unexpected error in repo.Get: %v", err)
			return ErrUnexpected
		}
	}

	if requester.ID == id || requester.Role == model.RoleAdmin {
		err = c.repo.Delete(id)
		if err != nil {
			re, ok := err.(*mock.ErrUserNotFound)
			if ok {
				c.l.Errorf("user with id %d not found", re.ID)
				return ErrUserNotFound
			} else {
				c.l.Errorf("controller.DeleteUser: unexpected error in repo.Delete: %v", err)
				return ErrUnexpected
			}
		}
	}
	return nil
}

func (c *User) RegeneratePfp(id int) error {
	m, err := c.repo.Get(id)
	if err != nil {
		re, ok := err.(*mock.ErrUserNotFound)
		if ok {
			c.l.Errorf("user with id %d not found", re.ID)
			return ErrUserNotFound
		} else {
			c.l.Errorf("controller.RegenerateLogo: unexpected error in repo.Get: %v", err)
			return ErrUnexpected
		}
	}

	m.Pfp, err = c.logo.GenerateLogo()
	if err != nil {
		c.l.Errorf("controller.RegenerateLogo: unexpected error in logo.GenerateLogo: %v", err)
		return ErrUnexpected
	}

	err = c.repo.Update(m)
	if err != nil {
		re, ok := err.(*mock.ErrUserNotFound)
		if ok {
			c.l.Errorf("user with id %d not found", re.ID)
			return ErrUserNotFound
		} else if err == mock.ErrUserUnapdatable {
			return ErrUnupdatableUser
		} else {
			c.l.Errorf("controller.RegenerateLogo: unexpected error in repo.Update: %v", err)
			return ErrUnexpected
		}
	}
	return nil
}

func (c *User) CheckCredentials(email, password string) (int, error) {
	m, err := c.repo.GetByEmail(email)
	if err != nil {
		_, ok := err.(*mock.ErrUserNotFound)
		if ok {
			c.l.Errorf("user with email %s not found", email)
			return -1, ErrUserNotFound
		} else {
			c.l.Errorf("controller.CheckCredentials: unexpected error in repo.GetByEmail: %v", err)
			return -1, ErrUnexpected
		}
	}

	if m.Password != password {
		return -1, ErrWrongPassword
	}

	return m.ID, nil
}
