package handler

import (
	"github.com/labstack/echo/v4"

	"github.com/RaivoKinne/Friends/utils"
	"github.com/RaivoKinne/Friends/web/templates"
)

func NotFound(c echo.Context) error {
	return utils.Render(c, templates.NotFound())
}
