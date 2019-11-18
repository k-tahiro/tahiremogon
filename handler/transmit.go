package handler

import (
	"net/http"
	"os/exec"
	"unsafe"

	"github.com/labstack/echo"

	myMw "github.com/k-tahiro/tahiremogon/middleware"
	"github.com/k-tahiro/tahiremogon/model"
	"github.com/k-tahiro/tahiremogon/util"
)

func TransmitCode(c echo.Context) error {
	cc := c.(*myMw.CustomContext)

	key := cc.Param("key")
	var code string

	cc.Session.Select("code").From("codes").Where("key = ?", key).Load(&code)
	if code == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Command Undefined")
	}

	cmd := "sudo /usr/local/bin/bto_ir_cmd -e -t" + " " + code
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	response := &model.TransmitResponse{
		Sucess: true,
		Label:  -1,
	}

	predictionModel := cc.PredictionModel
	if predictionModel != nil {
		label, _ := confirm(predictionModel)
		response.Label = label
	}

	return cc.JSON(http.StatusOK, response)
}

func confirm(predictionModel *myMw.PredictionModel) (int, error) {
	cmd := "sudo /usr/local/bin/camera.sh"
	filename, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return -1, err
	}

	input, err := util.ReadImage(*(*string)(unsafe.Pointer(&filename)))
	if err != nil {
		return -1, err
	}

	return predictionModel.Predict(input)
}
