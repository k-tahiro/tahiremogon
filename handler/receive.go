package handler

import (
	"net/http"
	"os/exec"

	"github.com/labstack/echo"

	myMw "github.com/k-tahiro/tahiremogon/middleware"
	"github.com/k-tahiro/tahiremogon/model"
)

func Receive(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	request := new(model.Request)
	if err := cc.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	cmd := "sudo /usr/local/bin/bto_ir_cmd -e -r | tail -n 1 | cut -f 2 -d : | cut -b 2- | tr -d '\n'"
	signal, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var command model.Command
	command.ID = request.ID
	command.Name = request.Name
	command.Signal = string(signal[:])

	sess := cc.Connection.NewSession(nil)
	sess.InsertInto("command").Columns("id", "name", "signal").Record(command).Exec()

	var response model.Response
	response.Success = true
	return cc.JSON(http.StatusOK, response)
}
