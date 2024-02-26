package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	"github.com/RaivoKinne/Friends/api/handler"
	"github.com/RaivoKinne/Friends/api/middleware"
)

func main() {
	router := echo.New()
	router.Static("/static", "web/static/")
	router.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	router.Debug = true
	router.RouteNotFound("/*", handler.NotFound)
	router.GET("/post", middleware.AuthMiddleware(handler.PostPageHandler))
	router.POST("/post/save", middleware.AuthMiddleware(handler.PostHandler))
	router.POST("/post/:id", middleware.AuthMiddleware(handler.DeletePostHandler))
	router.GET("/", handler.RegisterPageHandler)
	router.POST("/register", handler.RegisterHandler)
	router.GET("/login", handler.LoginPageHandler)
	router.POST("/login", handler.LoginHandler)
	router.GET("/logout", handler.LogoutHandler)
	router.Logger.Fatal(router.Start(":1323"))
}
