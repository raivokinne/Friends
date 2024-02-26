package handler

import (
	"database/sql"
	"net/mail"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/RaivoKinne/Friends/internal/database/model"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func FetchPosts(db *sql.DB) ([]model.Post, error) {
	posts, err := db.Query(`
        SELECT posts.id, posts.user_id, users.username, posts.content, posts.created_at
        FROM posts
        JOIN users ON posts.user_id = users.id
		ORDER BY posts.created_at DESC
    `)
	if err != nil {
		return nil, err
	}

	defer posts.Close()

	var postsList []model.Post

	for posts.Next() {
		var post model.Post
		var createdAt sql.NullString

		if err := posts.Scan(&post.ID, &post.UserID, &post.Username, &post.Content, &createdAt); err != nil {
			return nil, err
		}

		if createdAt.Valid {
			post.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt.String)
			if err != nil {
				return nil, err
			}
		} else {
			post.CreatedAt = time.Time{}
		}

		postsList = append(postsList, post)
	}

	return postsList, nil
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
