package handler

import (
	"log"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/RaivoKinne/Friends/internal/database"
	"github.com/RaivoKinne/Friends/internal/database/model"
	"github.com/RaivoKinne/Friends/utils"
	"github.com/RaivoKinne/Friends/web/templates/auth/login"
	"github.com/RaivoKinne/Friends/web/templates/auth/register"
)

var hPassword []byte

func RegisterPageHandler(c echo.Context) error {
	return utils.Render(c, register.Register(""))
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
		return utils.Render(c, register.Register("Please fill all fields"))
	}

	if len(user.Password) < 8 {
		return utils.Render(c, register.Register("Password must be at least 8 characters long"))
	}

	if !utils.IsValidEmail(user.Email) {
		return utils.Render(c, register.Register("Invalid email address"))
	}

	if user.Password != confirmPass {
		return utils.Render(c, register.Register("Passwords do not match"))
	}

	rows, err := db.Query("SELECT email FROM users WHERE email = ?", user.Email)
	if err != nil {
		log.Fatal("Error querying database:", err)
	}
	defer rows.Close()

	if rows.Next() {
		return utils.Render(c, register.Register("Email already exists"))
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

	return c.Redirect(302, "/posts")
}

func LoginPageHandler(c echo.Context) error {
	return utils.Render(c, login.Login(""))
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
		return utils.Render(c, login.Login("Invalid email or password"))
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

	return c.Redirect(302, "/posts")
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
