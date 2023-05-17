package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "ROLE_ADMIN"
	RoleUser  Role = "ROLE_USER"
)

func (r Role) String() string {
	return string(r)
}

type User struct {
	ID        uint      `json:"id" gorm:"column:id"`
	Username  string    `json:"username" gorm:"column:username;"`
	Email     string    `json:"email" gorm:"column:email;"`
	Password  string    `json:"-" gorm:"column:password;"`
	RolesAll  string    `json:"-" gorm:"column:roles;"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at"`
	Disabled  bool      `json:"-" gorm:"column:disabled;"`

	Roles    []string          `json:"-" gorm:"-"`
	RolesMap map[Role]struct{} `json:"-" gorm:"-"`
}

func (u *User) GetID() uint {
	return u.ID
}

func (u *User) ToModel() interface{} {
	return u
}

func (u *User) TableName() string {
	return "users"
}

// BeforeCreate sets RolesAll field in this u User.
func (u *User) BeforeCreate(_ *gorm.DB) error {
	u.RolesAll = strings.Join(u.Roles, " ")
	return nil
}

// AfterFind sets Roles and RolesMap from RolesAll.
func (u *User) AfterFind(_ *gorm.DB) error {
	roles := strings.Split(u.RolesAll, " ")
	u.Roles = u.Roles[:0]
	u.RolesMap = make(map[Role]struct{})
	for _, r := range roles {
		u.Roles = append(u.Roles, r)
		u.RolesMap[Role(r)] = struct{}{}
	}
	return nil
}

// Sanitize removes any private data.
func (u *User) Sanitize(_ map[string]struct{}) {
	u.Password = ""
}
