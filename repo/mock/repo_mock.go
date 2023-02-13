// The mock repo in this case will just be a map with the key being the ID of the user and
// the value being the user model.
package mock

import (
	"fmt"

	"github.com/vano2903/service-template/model"
	"github.com/vano2903/service-template/repo"
)

var (
	// This step sould be the first step when creating a repo component, check the struct with the interface.
	// To do so in golang you just write this: var _ Interface = new(Struct)
	// If this line of code doesn't compile then your struct doesn't implement the interface.
	_ repo.UserRepoer = new(RepoMock)

	//In this example we are using a custom error statically defined
	//and a custom error defined as a struct.

	//The static error is accesible from other packages and they can
	//check if the error returned by a function is the same as the static error.

	//The struct error is useful when we need to embed more informationl, for example the ID of the user
	ErrUserUnapdatable = fmt.Errorf("user can't be updated")
)

type ErrUserNotFound struct {
	ID      int
	Message string
}

func (e ErrUserNotFound) Error() string {
	return e.Message
}

type RepoMock struct {
	users  map[int]*model.User
	lastID int
}

// NewRepo returns a new mock repo.
func NewRepo() *RepoMock {
	return &RepoMock{
		users: make(map[int]*model.User),
	}
}

func (r *RepoMock) Create(u *model.User) (int, error) {
	r.lastID++
	u.ID = r.lastID
	r.users[u.ID] = u
	return int(r.lastID), nil
}

func (r *RepoMock) Get(id int) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		err := ErrUserNotFound{
			ID:      id,
			Message: fmt.Sprintf("user with id %d is not found", id),
		}
		return nil, err

	}
	return u, nil
}

func (r *RepoMock) GetByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	err := ErrUserNotFound{
		Message: fmt.Sprintf("user with email %s is not found", email),
	}
	return nil, err
}

func (r *RepoMock) Update(u *model.User) error {
	_, ok := r.users[u.ID]
	if !ok {
		err := ErrUserNotFound{
			ID:      u.ID,
			Message: fmt.Sprintf("user with id %d not found", u.ID),
		}
		return err
	}
	if r.users[u.ID].Role == model.RoleUnupdatable {
		return ErrUserUnapdatable
	}
	r.users[u.ID] = u
	return nil
}

func (r *RepoMock) Delete(id int) error {
	_, ok := r.users[id]
	if !ok {
		err := ErrUserNotFound{
			ID:      id,
			Message: fmt.Sprintf("user with id %d not found", id),
		}
		return err
	}
	delete(r.users, id)
	return nil
}

func (r *RepoMock) GetAll() []*model.User {
	users := make([]*model.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users
}
