package graph

import (
	"github.com/jinzhu/gorm"
	"twitter-clone/graph/model"
	"twitter-clone/middleware"
)

//go:generate go run github.com/99designs/gqlgen
type Resolver struct {
	User *model.User
	Post []*model.NewPost
	Db   *gorm.DB
	MW   middleware.RequireUser
}
