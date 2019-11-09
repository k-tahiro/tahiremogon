package handler

import (
	"net/http"

	"github.com/labstack/echo"

	myMw "github.com/k-tahiro/tahiremogon/middleware"
	"github.com/k-tahiro/tahiremogon/model"
)

func Commands(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	var commands []model.CommandJSON
	sess := cc.Connection.NewSession(nil)
	sess.Select("id", "name", "signal").From("command").Load(&commands)

	return cc.JSON(http.StatusOK, commands)
}
