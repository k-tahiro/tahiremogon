package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"

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

	filename, err := exec.Command("/usr/local/bin/transmit.sh", signal).Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Command Execution Failed")
	}

	// TODO: 画像認識処理
	fmt.Println(filename)

	var response model.Response
	response.Success = true
	return cc.JSON(http.StatusOK, response)
}

func Receive(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	request := new(model.Request)
	if err := cc.Bind(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad Request Body")
	}

	signal, err := exec.Command("/usr/local/bin/receive.sh").Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Command Execution Failed")
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

func ReceiveImage(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	filename := os.Getenv("IMG_DIR") + "/image_" + time.Now().String() + ".jpg"

	//-----------
	// Read file
	//-----------

	// Source
	file, err := cc.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	fmt.Println(filename)

	var response model.Response
	response.Success = true
	return cc.JSON(http.StatusOK, response)
}
