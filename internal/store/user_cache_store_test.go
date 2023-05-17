package store

import (
	"context"
	"errors"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/zacscoding/go-rest-template/internal/model"
)

func (s *CacheStoreSuite) TestUserStore_Save() {
	user := model.User{
		ID:       1,
		Username: "user1",
		Email:    "user1@email.com",
		Password: "userpassword",
		Roles:    []string{string(model.RoleUser)},
	}
	s.userStoreMock.On("Save", mock.Anything, &user).Return(nil)

	err := s.userStore.Save(context.TODO(), &user)

	s.NoError(err)
	s.userStoreMock.AssertCalled(s.T(), "Save", mock.Anything, &user)
}

func (s *CacheStoreSuite) TestUserStore_Save_Fail() {
	user := model.User{
		ID:       1,
		Username: "user1",
		Email:    "user1@email.com",
		Password: "userpassword",
		Roles:    []string{string(model.RoleUser)},
	}
	s.userStoreMock.On("Save", mock.Anything, &user).Return(errors.New("force err"))

	err := s.userStore.Save(context.TODO(), &user)

	s.Error(err)
	s.Contains(err.Error(), "force err")
}

func (s *CacheStoreSuite) TestUserStore_FindByEmail_NoCacheHit() {
	user := model.User{
		ID:        1,
		Username:  "user1",
		Email:     "user@email.com",
		Password:  "userpass",
		RolesAll:  string(model.RoleUser),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Add(time.Second),
		Disabled:  false,
		Roles:     []string{string(model.RoleUser)},
		RolesMap:  map[model.Role]struct{}{model.RoleUser: {}},
	}
	s.mpMock.On("RecordCache", mock.Anything, mock.Anything)
	s.userStoreMock.On("FindByEmail", mock.Anything, user.Email).Return(&user, nil)

	find, err := s.userStore.FindByEmail(context.TODO(), user.Email)

	s.NoError(err)
	s.checkUser(&user, find)
	s.userStoreMock.AssertCalled(s.T(), "FindByEmail", mock.Anything, user.Email)
	s.userStoreMock.AssertNumberOfCalls(s.T(), "FindByEmail", 1)
	s.mpMock.AssertCalled(s.T(), "RecordCache", mock.Anything, false)
}

func (s *CacheStoreSuite) TestUserStore_FindByEmail_CacheHit() {
	user := model.User{
		ID:        1,
		Username:  "user1",
		Email:     "user@email.com",
		Password:  "userpass",
		RolesAll:  string(model.RoleUser),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now().Add(time.Second),
		Disabled:  false,
		Roles:     []string{string(model.RoleUser)},
		RolesMap:  map[model.Role]struct{}{model.RoleUser: {}},
	}
	s.mpMock.On("RecordCache", mock.Anything, mock.Anything)
	s.userStoreMock.On("FindByEmail", mock.Anything, user.Email).Return(&user, nil)

	_, err := s.userStore.FindByEmail(context.TODO(), user.Email)
	s.NoError(err)
	find, err := s.userStore.FindByEmail(context.TODO(), user.Email)

	s.NoError(err)
	s.checkUser(&user, find)
	s.userStoreMock.AssertCalled(s.T(), "FindByEmail", mock.Anything, user.Email)
	s.userStoreMock.AssertNumberOfCalls(s.T(), "FindByEmail", 1)
	s.mpMock.AssertCalled(s.T(), "RecordCache", mock.Anything, true)
}

func (s *CacheStoreSuite) checkUser(expected, actual *model.User) {
	s.Equal(expected.ID, actual.ID)
	s.Equal(expected.Username, actual.Username)
	s.Equal(expected.Email, actual.Email)
	s.Equal(expected.Password, actual.Password)
	s.Equal(expected.RolesAll, actual.RolesAll)
	s.WithinDuration(expected.CreatedAt, actual.CreatedAt, time.Second)
	s.WithinDuration(expected.UpdatedAt, actual.UpdatedAt, time.Second)
	s.Equal(expected.Disabled, actual.Disabled)
	s.Equal(expected.Roles, actual.Roles)
	s.Equal(expected.RolesMap, actual.RolesMap)
}
