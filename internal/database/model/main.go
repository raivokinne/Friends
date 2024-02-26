package model

import (
	"time"
)

type User struct {
	ID       int
	Username string
	Password string
	Email    string
}

func (u *User) TableName() string {
	return "users"
}

type Post struct {
	ID        int
	UserID    int
	Username  string
	Content   string
	CreatedAt time.Time
}

func (p *Post) TableName() string {
	return "posts"
}
