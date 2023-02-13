package repo

import (
	"github.com/vano2903/service-template/model"
)

// Interfaces ends with -er
type (
	//This interface has the methods declarations for the
	//repo components.
	UserRepoer interface {
		Create(u *model.User) (id int, err error)
		Get(id int) (*model.User, error)
		GetByEmail(email string) (*model.User, error)
		Update(u *model.User) error
		Delete(id int) error
		GetAll() []*model.User
	}
)
