package controller

//This file has the interface for the controller.
//This can be useful to remember the methods to implement but it is not necessary
//because you shouldn't implement different version of the controller as it is the business logic.
//It's either a different version of the service or you should not implement a different version for the controller
//because it is the business logic. (for example there shouldn't be a mock for the controller)

import (
	"github.com/vano2903/service-template/model"
)

type (
	UserControllerer interface {
		CreateUser(firstName, lastName, email, password, role string) (int, error)
		GetUser(id int) (*model.User, error)
		GetAllUsers() []*model.User
		UpdateUser(requesterId int, u *model.User) error
		DeleteUser(requesterId int, id int) error
		RegeneratePfp(id int) error
		CheckCredentials(email, password string) (int, error)
	}
)
