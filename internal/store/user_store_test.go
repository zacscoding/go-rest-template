package store

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zacscoding/go-rest-template/internal/model"
	"github.com/zacscoding/go-rest-template/pkg/database"
)

func (s *StoreSuite) TestSave() {
	u := model.User{
		Username: "user1",
		Email:    "user1@email.com",
		Password: "user1pass",
		Roles:    []string{string(model.RoleUser)},
	}

	err := s.userStore.Save(context.TODO(), &u)

	s.NoError(err)
	s.NotEqualValues(0, u.ID)
	find, err := s.userStore.FindByEmail(context.TODO(), u.Email)
	s.NoError(err)
	s.Equal(u.ID, find.ID)
	s.Equal(u.Username, find.Username)
	s.Equal(u.Email, find.Email)
	s.Equal(u.Password, find.Password)
	s.WithinDuration(find.CreatedAt, time.Now(), time.Second)
	s.WithinDuration(find.UpdatedAt, time.Now(), time.Second)
}

func (s *StoreSuite) TestSave_Fail() {
	saved := model.User{Email: "user1@email.com", Roles: []string{string(model.RoleUser)}}
	s.NoError(s.userStore.Save(context.TODO(), &saved))

	s.T().Run("Duplicate Email", func(t *testing.T) {
		err := s.userStore.Save(context.TODO(), &model.User{
			Username: saved.Username + "_1",
			Email:    saved.Email,
			Password: saved.Password + "_1",
			Roles:    saved.Roles,
		})

		assert.Error(t, err)
		assert.Equal(t, database.ErrKeyConflict, err)
	})
}

func (s *StoreSuite) TestFindByEmail() {
	saved := model.User{Email: "user1@email.com", Roles: []string{string(model.RoleUser)}}
	s.NoError(s.userStore.Save(context.TODO(), &saved))

	find, err := s.userStore.FindByEmail(context.TODO(), saved.Email)
	s.NoError(err)
	s.Equal(saved.ID, find.ID)
	s.Equal(saved.Username, find.Username)
	s.Equal(saved.Email, find.Email)
	s.Equal(saved.Password, find.Password)
	s.WithinDuration(find.CreatedAt, time.Now(), time.Second)
	s.WithinDuration(find.UpdatedAt, time.Now(), time.Second)
	s.Contains(find.Roles, string(model.RoleUser))
	s.Contains(find.RolesMap, model.RoleUser)
}

func (s *StoreSuite) TestFindByEmail_Fail() {
	saved := model.User{Email: "user1@email.com", Roles: []string{string(model.RoleUser)}}
	s.NoError(s.userStore.Save(context.TODO(), &saved))

	s.T().Run("Not Found", func(t *testing.T) {
		find, err := s.userStore.FindByEmail(context.TODO(), "user2@email.com")

		assert.Nil(t, find)
		assert.Error(t, err)
		assert.Equal(t, database.ErrRecordNotFound, err)
	})

	s.T().Run("Not Found With Empty Email", func(t *testing.T) {
		find, err := s.userStore.FindByEmail(context.TODO(), "")

		assert.Nil(t, find)
		assert.Error(t, err)
		assert.Equal(t, database.ErrRecordNotFound, err)
	})
}
