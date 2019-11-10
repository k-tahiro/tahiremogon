package handler

import (
	"net/http"
	"os/exec"
	"unsafe"

	"github.com/labstack/echo"

	myMw "github.com/k-tahiro/tahiremogon/middleware"
	"github.com/k-tahiro/tahiremogon/model"
)

func CreateCode(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	key := cc.Param("key")
	cmd := "sudo /usr/local/bin/bto_ir_cmd -e -r | tail -n 1 | cut -f 2 -d : | cut -b 2- | tr -d '\n'"
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	code := *(*string)(unsafe.Pointer(&output))
	sCode := &model.Code{
		Key:  key,
		Code: code,
	}

	sess := cc.Connection.NewSession(nil)
	sess.InsertInto("command").Columns("key", "code").Record(sCode).Exec()

	return cc.JSON(http.StatusOK, sCode)
}

func ReadCodes(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	var codes []model.Code
	sess := cc.Connection.NewSession(nil)
	sess.Select("id", "key", "code").From("codes").Load(&codes)

	return cc.JSON(http.StatusOK, codes)
}

func DeleteCode(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	key := cc.Param("key")
	var code model.Code

	sess := cc.Connection.NewSession(nil)
	sess.Select("*").From("codes").Where("key = ?", key).Load(&code)
	sess.DeleteFrom("codes").Where("key = ?", key).Exec()

	return cc.JSON(http.StatusOK, code)
}
