package handler

import (
	"net/http"
	"os/exec"

	"github.com/labstack/echo"

	myMw "../middleware"
	"../model"
)

// Handler
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func Commands(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	var commands []model.Command
	sess := cc.Connection.NewSession(nil)
	sess.Select("id","name").From("command").Load(&commands)
	
    return cc.JSON(http.StatusOK, commands)
}

func Transmit(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	id := cc.Param("id")
	var signal string
	sess := cc.Connection.NewSession(nil)
	sess.Select("signal").From("command").Where("id = ?", id).Load(&signal)

	err := exec.Command("/usr/local/bin/bto_ir_cmd", "-e", "-t", signal).Run()
	var response model.Response
	if err == nil {
		response.Success = true
	} else {
		response.Success = false
	}

	return cc.JSON(http.StatusOK, response)
}

func Receive(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	request := new(model.Request)
	if err := cc.Bind(request); err != nil {
		return err
	}

	out, err := exec.Command("/usr/local/bin/bto_ir_cmd", "-e", "-r").Output()
	return cc.JSON(http.StatusOK, out)
}
