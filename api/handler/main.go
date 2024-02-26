package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/RaivoKinne/Friends/internal/database"
	"github.com/RaivoKinne/Friends/internal/database/model"
	"github.com/RaivoKinne/Friends/web/templates"
	"github.com/RaivoKinne/Friends/web/templates/auth/login"
	"github.com/RaivoKinne/Friends/web/templates/auth/register"
	"github.com/RaivoKinne/Friends/web/templates/post"
)

var hPassword []byte

func NotFound(c echo.Context) error {
	return Render(c, templates.NotFound())
}

func PostPageHandler(c echo.Context) error {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	posts, err := FetchPosts(db)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	return Render(c, post.Index(posts))
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

	return c.Redirect(302, "/post")
}

func RegisterPageHandler(c echo.Context) error {
	return Render(c, register.Register(""))
}

func RegisterHandler(c echo.Context) error {
	user := &model.User{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
		Email:    c.FormValue("email"),
	}

	db, err := database.Connect()
	if err != nil {
		log.Println("Error connecting to the database:", err)
		return err
	}
	defer db.Close()

	confirmPass := c.FormValue("confirmPassword")

	if user.Username == "" || user.Password == "" || user.Email == "" {
		return Render(c, register.Register("Please fill all fields"))
	}

	if len(user.Password) < 8 {
		return Render(c, register.Register("Password must be at least 8 characters long"))
	}

	if !IsValidEmail(user.Email) {
		return Render(c, register.Register("Invalid email address"))
	}

	if user.Password != confirmPass {
		return Render(c, register.Register("Passwords do not match"))
	}

	rows, err := db.Query("SELECT email FROM users WHERE email = ?", user.Email)
	if err != nil {
		log.Fatal("Error querying database:", err)
	}
	defer rows.Close()

	if rows.Next() {
		return Render(c, register.Register("Email already exists"))
	} else if err := rows.Err(); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return err
	}

	hPassword = hashedPassword

	stmt, err := db.Prepare("INSERT INTO users (username, password, email) VALUES (?, ?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Username, hashedPassword, user.Email)
	if err != nil {
		log.Println("Error inserting user into the database:", err)
		return err
	}

	row := db.QueryRow("SELECT id FROM users WHERE username = ?", user.Username)
	err = row.Scan(&user.ID)
	if err != nil {
		log.Println("Error querying database:", err)
		return err
	}

	sess, err := session.Get("session", c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	sess.Values["Authenticated"] = true
	sess.Values["UserId"] = user.ID

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		log.Println("Error saving session:", err)
		return err
	}

	return c.Redirect(302, "/post")
}

func LoginPageHandler(c echo.Context) error {
	return Render(c, login.Login(""))
}

func LoginHandler(c echo.Context) error {
	user := &model.User{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	db, err := database.Connect()
	if err != nil {
		log.Println("Error connecting to the database:", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT password FROM users WHERE email = ?")
	if err != nil {
		log.Println("Error preparing statement:", err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(user.Email).Scan(&hPassword)
	if err != nil {
		log.Println("Error querying database:", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hPassword), []byte(user.Password))
	if err != nil {
		return Render(c, login.Login("Invalid email or password"))
	}

	row := db.QueryRow("SELECT id FROM users WHERE email = ?", user.Email)
	err = row.Scan(&user.ID)
	if err != nil {
		log.Println("Error querying database:", err)
	}

	sess, err := session.Get("session", c)
	if err != nil {
		log.Println("Error getting session:", err)
	}

	sess.Values["Authenticated"] = true
	sess.Values["UserId"] = user.ID

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		log.Println("Error saving session:", err)
	}

	return c.Redirect(302, "/post")
}

func LogoutHandler(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		log.Println("Error getting session:", err)
		return err
	}

	sess.Values["UserId"] = 0
	sess.Values["Authenticated"] = false
	sess.Save(c.Request(), c.Response())

	return c.Redirect(302, "/")
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

	return c.Redirect(http.StatusFound, "/post")
}
