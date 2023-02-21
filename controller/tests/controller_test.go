package controller

import (
	"log"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/vano2903/service-template/config"
	"github.com/vano2903/service-template/controller"
	"github.com/vano2903/service-template/model"
	"github.com/vano2903/service-template/pkg/logger"
	"github.com/vano2903/service-template/providers/logo"
	"github.com/vano2903/service-template/repo/mock"
	"gotest.tools/v3/assert"
)

var (
	c *controller.User
	r *mock.RepoMock
	l *logrus.Logger
)

func init() {
	conf := config.Config{}
	conf.Log.Level = "debug"
	conf.Log.Type = "text"
	conf.Database.Driver = "mock"
	l = logger.NewLogger(conf.Log.Level, conf.Log.Type)

	if conf.Database.Driver != "mock" {
		log.Fatal("only mock database is supported in this example")
	}

	r = mock.NewRepo()
	logoService := logo.NewServiceLogo(conf.Services.Logo.ApiKey, conf.Services.Logo.BaseUrl)

	c = controller.NewUserController(r, logoService, l)

	// GenerateExampleEntries(l, c)
}

// create users, cases:
// [x] create a user
// [x] block creating a user with the same email
func TestCreateUser(t *testing.T) {
	firstName := "name"
	lastName := "lastname"
	email := "amazingemail@gmail.com"
	password := "password"
	var id int
	t.Run("create user", func(t *testing.T) {

		id, err := c.CreateUser(firstName, lastName, email, password, model.RoleUser)
		if err != nil {
			t.Errorf("unable to create user: %v", err)
		}

		//we get the user from the repo directly to make sure
		//it's actually been created correctly
		user, err := r.Get(id)
		if err != nil {
			re, ok := err.(*mock.ErrUserNotFound)
			if ok {
				t.Errorf("user was not created correctly: %v", re)
			} else {
				t.Errorf("unexpected error from repo: %v", err)
			}
		}

		assert.Equal(t, user.FirstName, firstName)
		assert.Equal(t, user.LastName, lastName)
		assert.Equal(t, user.Email, email)
		assert.Equal(t, user.Password, password)
		assert.Equal(t, user.Role, model.RoleUser)
	})

	t.Run("create duplicate user", func(t *testing.T) {
		var err error
		id, err = c.CreateUser(firstName, lastName, email, password, model.RoleUser)
		if err == nil {
			t.Error("duplicate user was created, should not happen")
		}

	})

	t.Cleanup(func() {
		if err := r.Delete(id); err != nil {
			t.Errorf("error cleaning up test case: %v", err)
		}
	})
}

// update user, cases:
// [x] update a roleUser user
// [x] unable to update a roleUnupdatable user
// [x] block roleUser user to update another user
// [x] allow roleAdmin user to update any user
// [x] block roleAdmin from updating an roleUnupdatable user
func TestUpdateUser(t *testing.T) {
	firstName := "name"
	lastName := "lastname"
	email := "test@test.com"
	password := "password"
	t.Run("update roleUser user", func(t *testing.T) {
		id, err := c.CreateUser(firstName, lastName, email, password, model.RoleUser)
		if err != nil {
			t.Errorf("unable to create user: %v", err)
		}

		updatedUser, err := r.Get(id)
		if err != nil {
			t.Errorf("unable to get user from repo: %v", err)
		}
		t.Logf("user created: %v", updatedUser)

		updatedFirstName := "updatedName"
		updatedPassword := "updatedPassword"
		updatedUser.FirstName = updatedFirstName
		updatedUser.Password = updatedPassword

		if err := c.UpdateUser(id, updatedUser); err != nil {
			t.Errorf("unable to update user: %v", err)
		}

		user, err := r.Get(id)
		t.Log(user)
		if err != nil {
			t.Error("unable to get user from repo")
		}

		assert.Equal(t, user.FirstName, updatedFirstName)
		assert.Equal(t, user.Password, updatedPassword)
		t.Cleanup(func() {
			if err := r.Delete(id); err != nil {
				t.Errorf("error cleaning up test case: %v", err)
			}
		})
	})

	t.Run("unable to update roleUnupdatable user", func(t *testing.T) {
		id, err := c.CreateUser(firstName, lastName, email, password, model.RoleUnupdatable)
		if err != nil {
			t.Errorf("unable to create user: %v", err)
		}

		updatedFirstName := "updatedName"
		updatedPassword := "updatedPassword"
		updatedUser := model.User{
			FirstName: updatedFirstName,
			LastName:  lastName,
			Email:     email,
			Password:  updatedPassword,
			Role:      model.RoleUser,
		}
		if err := c.UpdateUser(id, &updatedUser); err == nil {
			t.Error("unupdatable user was updated, should not happen")
		}
		t.Cleanup(func() {
			if err := r.Delete(id); err != nil {
				t.Errorf("error cleaning up test case: %v", err)
			}
		})
	})

	t.Run("block roleUser user to update another user", func(t *testing.T) {
		idUpdater, err := c.CreateUser(firstName, lastName, email, password, model.RoleUser)
		if err != nil {
			t.Errorf("unable to create user: %v", err)
		}

		toUpdate := model.User{
			FirstName: "toUpdate",
			LastName:  "toUpdate",
			Email:     "test2@test.com",
			Password:  "toUpdate",
			Role:      model.RoleUser,
		}

		idUpdated, err := c.CreateUser(toUpdate.FirstName, toUpdate.LastName, toUpdate.Email, toUpdate.Password, toUpdate.Role)
		if err != nil {
			t.Errorf("unable to create user: %v", err)
		}

		if err := c.UpdateUser(idUpdater, &toUpdate); err == nil {
			t.Error("roleUser user was able to update another user, should not happen")
		}

		t.Cleanup(func() {
			if err := r.Delete(idUpdater); err != nil {
				t.Errorf("error cleaning up test case: %v", err)
			}
			if err := r.Delete(idUpdated); err != nil {
				t.Errorf("error cleaning up test case: %v", err)
			}
		})
	})

	t.Run("allow roleAdmin user to update any user", func(t *testing.T) {
		idUpdater, err := c.CreateUser(firstName, lastName, email, password, model.RoleAdmin)
		if err != nil {
			t.Errorf("unable to create user: %v", err)
		}

		toUpdate := &model.User{
			FirstName: "toUpdate",
			LastName:  "toUpdate",
			Email:     "test2@test.com",
			Password:  "toUpdate",
			Role:      model.RoleUser,
		}

		idUpdated, err := c.CreateUser(toUpdate.FirstName, toUpdate.LastName, toUpdate.Email, toUpdate.Password, toUpdate.Role)
		if err != nil {
			t.Errorf("unable to create user: %v", err)
		}

		toUpdate, err = r.Get(idUpdated)
		if err != nil {
			t.Errorf("unable to get user from repo: %v", err)
		}

		if err := c.UpdateUser(idUpdater, toUpdate); err != nil {
			t.Errorf("roleAdmin user was not able to update another user: %v", err)
		}

		t.Cleanup(func() {
			if err := r.Delete(idUpdater); err != nil {
				t.Errorf("error cleaning up test case: %v", err)
			}
			if err := r.Delete(idUpdated); err != nil {
				t.Errorf("error cleaning up test case: %v", err)
			}
		})
	})

	t.Run("block roleAdmin user from updating roleUnupdatable user", func(t *testing.T) {
		idUpdater, err := c.CreateUser(firstName, lastName, email, password, model.RoleAdmin)
		if err != nil {
			t.Errorf("unable to create user: %v", err)
		}

		toUpdate := model.User{
			FirstName: "toUpdate",
			LastName:  "toUpdate",
			Email:     "test2@test.com",
			Password:  "toUpdate",
			Role:      model.RoleUnupdatable,
		}

		idUpdated, err := c.CreateUser(toUpdate.FirstName, toUpdate.LastName, toUpdate.Email, toUpdate.Password, toUpdate.Role)
		if err != nil {
			t.Errorf("unable to create user: %v", err)
		}

		if err := c.UpdateUser(idUpdater, &toUpdate); err == nil {
			t.Error("roleAdmin user was able to update roleUnupdatable user, should not happen")
		}

		t.Cleanup(func() {
			if err := r.Delete(idUpdater); err != nil {
				t.Errorf("error cleaning up test case: %v", err)
			}
			if err := r.Delete(idUpdated); err != nil {
				t.Errorf("error cleaning up test case: %v", err)
			}
		})
	})
}

// regenerate user's logo
func TestRegenerateLogo(t *testing.T) {
	firstName := "name"
	lastName := "lastname"
	email := "test@test.com"
	password := "password"

	id, err := c.CreateUser(firstName, lastName, email, password, model.RoleUser)
	if err != nil {
		t.Errorf("unable to create user: %v", err)
	}

	u, err := r.Get(id)
	if err != nil {
		t.Errorf("unable to get user from repo: %v", err)
	}
	prevLogo := u.Pfp

	if err := c.RegeneratePfp(id); err != nil {
		t.Errorf("unable to regenerate logo: %v", err)
	}

	u, err = r.Get(id)
	if err != nil {
		t.Errorf("unable to get user from repo: %v", err)
	}

	if u.Pfp == prevLogo {
		t.Error("logo was not regenerated")
	}

	t.Cleanup(func() {
		if err := r.Delete(id); err != nil {
			t.Errorf("error cleaning up test case: %v", err)
		}
	})
}
