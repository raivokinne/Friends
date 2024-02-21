package middleware

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return err
		}
		if sess.Values["Authenticated"] != true {
			return c.Redirect(http.StatusFound, "/")
		}
		return next(c)
	}
}
