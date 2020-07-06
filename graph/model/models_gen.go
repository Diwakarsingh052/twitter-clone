package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Follow struct {
	MyEmail       string `json:"MyEmail" gorm:"primary_key"`
	FollowedEmail string `json:"followedEmail"  gorm:"primary_key"`
}

type FollowUser struct {
	Email string `json:"Email"`
}

type FollowedUser struct {
	Email string `json:"email"`
}

type NewPost struct {
	Text string `json:"text"`
}

type Post struct {
	Email string    `json:"email"`
	Text  string    `json:"text"`
	Time  time.Time `json:"time"`
}

type User struct {
	gorm.Model
	Email        string  `json:"email" gorm:"not null;unique_index"`
	Password     string  `json:"password" gorm:"-"`
	PasswordHash *string `json:"passwordHash"`
	Remember     *string `json:"remember" gorm:"-"`
	RememberHash *string `json:"rememberHash"`
}
type UserToken struct {
	Token string
}
