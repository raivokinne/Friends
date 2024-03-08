package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/RaivoKinne/Friends/internal/database"
	"github.com/RaivoKinne/Friends/internal/database/model"
	"github.com/RaivoKinne/Friends/utils"
	"github.com/RaivoKinne/Friends/web/templates/post"
)

func PostPageHandler(c echo.Context) error {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	posts, err := utils.FetchPosts(db)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	sess, err := session.Get("session", c)

	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	userId, ok := sess.Values["UserId"].(int)
	if !ok {
		log.Println("User ID is nil or not of type int")
	}

	return utils.Render(c, post.Index(posts, userId))
}

func PostHandler(c echo.Context) error {
	sess, err := session.Get("session", c)

	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	if auth, ok := sess.Values["Authenticated"].(bool); !ok || !auth {
		log.Println("User is not authenticated")
		return c.Redirect(http.StatusFound, "/")
	}

	userId, ok := sess.Values["UserId"].(int)
	if !ok {
		log.Println("User ID is nil or not of type int")
		return c.Redirect(http.StatusFound, "/")
	}

	post := &model.Post{
		UserID:    userId,
		Content:   c.FormValue("content"),
		CreatedAt: time.Now().UTC(),
	}

	db, err := database.Connect()
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id INT,
		content VARCHAR(255),
		created_at DATETIME,
        FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		log.Println("Error creating messages table:", err)
		return err
	}

	stmt, err := db.Prepare("INSERT INTO posts (user_id, content, created_at) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(post.UserID, post.Content, post.CreatedAt)
	if err != nil {
		log.Println("Error inserting message into the database:", err)
		return err
	}

	return c.Redirect(302, "/posts")
}

func DeletePostHandler(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	if auth, ok := sess.Values["Authenticated"].(bool); !ok || !auth {
		log.Println("User is not authenticated")
		return c.Redirect(http.StatusFound, "/")
	}

	userId, ok := sess.Values["UserId"].(int)

	if !ok {
		log.Println("User ID is nil or not of type int")
		return c.Redirect(http.StatusFound, "/")
	}

	db, err := database.Connect()

	if err != nil {
		log.Println("Error connecting to the database:", err)
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM posts WHERE id = ? AND user_id = ?")

	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(c.Param("id"), userId)

	if err != nil {
		log.Println("Error deleting post:", err)
		return err
	}

	return c.Redirect(http.StatusFound, "/posts")
}
