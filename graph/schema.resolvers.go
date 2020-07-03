package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"log"

	"time"
	"twitter-clone/graph/generated"
	"twitter-clone/graph/model"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*model.Post, error) {
	var token = r.MW.Us.Token
	user, err := r.MW.Us.ByRemember(token)
	if err != nil {
		log.Fatal(err)
	}
	p := &model.Post{
		Email: user.Email,
		Text:  input.Text,
		Time:  time.Now().UTC(),
	}

	if err := r.Db.AutoMigrate(&model.Post{}).Error; err != nil {
		return nil, err
	}
	err = r.Db.Create(&p).Error
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (r *mutationResolver) FollowUser(ctx context.Context, input model.FollowUser) (*model.Follow, error) {
	var token = r.MW.Us.Token
	user, err := r.MW.Us.ByRemember(token)
	if err != nil {
		log.Fatal(err)
	}
	if user.Email == input.Email {
		return nil, errors.New("you cannot follow yourself")
	}
	followUser, err := r.MW.Us.ByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	u := &model.Follow{
		MyEmail:       user.Email,
		FollowedEmail: followUser.Email,
	}

	if err := r.Db.AutoMigrate(&model.Follow{}).Error; err != nil {
		return nil, err
	}
	err = r.Db.Create(&u).Error
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *queryResolver) AllUsers(ctx context.Context) ([]*model.User, error) {
	var user []model.User
	err := r.Db.Find(&user).Error
	fmt.Println(user)
	var users = make([]*model.User, len(user))

	for i := 0; i < len(user); i++ {
		users[i] = &user[i]
	}
	if err != nil {

		return nil, err
	}

	return users, nil
}

func (r *queryResolver) MyPost(ctx context.Context) ([]*model.Post, error) {
	var post []model.Post
	var token = r.MW.Us.Token
	user, err := r.MW.Us.ByRemember(token)
	if err != nil {
		log.Fatal(err)
	}
	err = r.Db.Order("time desc").Where("email = ?", user.Email).Find(&post).Error

	var posts = make([]*model.Post, len(post))
	for i := 0; i < len(post); i++ {
		posts[i] = &post[i]

	}

	if err != nil {

		return nil, err
	}

	return posts, nil
}

func (r *queryResolver) FollowedPost(ctx context.Context) ([]*model.Post, error) {
	var post []model.Post
	var follow []model.Follow
	var token = r.MW.Us.Token
	user, err := r.MW.Us.ByRemember(token)
	if err != nil {
		log.Fatal(err)
	}

	err = r.Db.Where("my_email = ?", user.Email).Find(&follow).Error

	if err != nil {
		return nil, err
	}

	var posts []model.Post
	for i := 0; i < len(follow); i++ {
		err = r.Db.Order("time asc").Where("email = ?", follow[i].FollowedEmail).Find(&post).Error
		if err != nil {
			return nil, err
		}
		posts = append(posts, post...)
		fmt.Println(post)
	}


	Sort(posts)

	var finalPosts = make([]*model.Post, len(posts))
	for i := 0; i < len(posts); i++ {
		finalPosts[i] = &posts[i]

	}
	if err != nil {

		return nil, err
	}

	return finalPosts, nil
}

func (r *queryResolver) FollowedUser(ctx context.Context) ([]*model.FollowedUser, error) {
	var follow []model.Follow

	var token = r.MW.Us.Token
	user, err := r.MW.Us.ByRemember(token)
	if err != nil {
		log.Fatal(err)
	}

	err = r.Db.Where("my_email = ?", user.Email).Find(&follow).Error
	if err != nil {
		return nil, err
	}

	var followedTemp []model.FollowedUser = make([]model.FollowedUser, len(follow))

	for i, f := range follow {
		followedTemp[i].Email = f.FollowedEmail
	}

	var followed = make([]*model.FollowedUser, len(followedTemp))
	for i := 0; i < len(followedTemp); i++ {
		followed[i] = &followedTemp[i]

	}

	return followed, nil
}
func Sort(posts []model.Post) {
	len := len(posts)
	for i := 0; i < len-1; i++ {
		for j := 0; j < len-i-1; j++ {
			if posts[j].Time.Before(posts[j+1].Time) {
				posts[j], posts[j+1] = posts[j+1], posts[j]
			}
		}
	}

}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) Allusers(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	panic(fmt.Errorf("not implemented"))
}
