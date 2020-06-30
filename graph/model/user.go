package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
	"os"

	"twitter-clone/hash"
	"twitter-clone/rand"
)

var (
	// ErrNotFound is returned when a resource cannot be found
	// in the database.
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided
	// to a method like Delete.
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password
	// is used when attempting to authenticate a user.
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

const userPwPepper = "secret-random-string"
const hmacSecretKey = "secret-hmac-key"

func NewUserService(connectionInfo string) (*UserService, error) {
	Db, err := gorm.Open("mysql", connectionInfo)
	if err != nil {
		return nil, err
	}
	Db.Exec("Use " + os.Getenv("DATABASE"))
	//Db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &UserService{
		DB:   Db,
		hmac: hmac,
	}, nil
}

type UserService struct {
	DB    *gorm.DB
	hmac  hash.HMAC
	Token string
}

func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.DB.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := us.hmac.Hash(token)
	err := first(us.DB.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	var Hash string
	if foundUser.PasswordHash != nil {
		Hash = *foundUser.PasswordHash
	}

	err = bcrypt.CompareHashAndPassword([]byte(Hash), []byte(password+userPwPepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}

	return foundUser, nil
}
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper)

	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	var pass string = string(hashedBytes)
	if pass != "" {
		user.PasswordHash = &pass
	}

	user.Password = ""

	if user.Remember == nil {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}

		user.Remember = &token
	}

	var has = us.hmac.Hash(*user.Remember)
	user.RememberHash = &has
	fmt.Println(*user.Remember)
	return us.DB.Create(user).Error
}

func (us *UserService) Close() error {
	return us.DB.Close()
}
func (us *UserService) AutoMigrate() error {
	//us.DB.DropTableIfExists(&User{})
	if err := us.DB.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}
func (us *UserService) Update(user *User) error {
	if *user.Remember != "" {
		var rem = us.hmac.Hash(*user.Remember)
		user.RememberHash = &rem
	}

	return us.DB.Model(user).Update("remember_hash", user.RememberHash).Error

}
func (us *UserService) SetToken(token string) *UserToken {
	return &UserToken{Token: token}
}
