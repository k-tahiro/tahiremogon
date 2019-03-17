package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"

	myMw "../middleware"
	"../model"
)

func Commands(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	var commands []model.CommandJSON
	sess := cc.Connection.NewSession(nil)
	sess.Select("id", "name", "signal").From("command").Load(&commands)

	return cc.JSON(http.StatusOK, commands)
}

func Transmit(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	id := cc.Param("id")
	var signal string
	sess := cc.Connection.NewSession(nil)
	sess.Select("signal").From("command").Where("id = ?", id).Load(&signal)
	if signal == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Command Undefined")
	}

	cmd := "sudo /usr/local/bin/bto_ir_cmd -e -t" + " " + signal
	if output, err := execCommand(os.Getenv("MODE"), cmd, cc); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var response model.Response
	response.Success = true
	return cc.JSON(http.StatusOK, response)
}

func Receive(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	request := new(model.Request)
	if err := cc.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	cmd := "sudo /usr/local/bin/bto_ir_cmd -e -r | tail -n 1 | cut -f 2 -d : | cut -b 2- | tr -d '\n'"
	if signal, err := execCommand(os.Getenv("MODE"), cmd, cc); err != nil {
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
