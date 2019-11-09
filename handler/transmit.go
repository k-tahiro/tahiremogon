package handler

import (
	"errors"
	"fmt"
	"image"
	"net/http"
	"os"
	"os/exec"
	"unsafe"

	"github.com/labstack/echo"
	"github.com/nfnt/resize"
	"gorgonia.org/tensor"

	myMw "github.com/k-tahiro/tahiremogon/middleware"
	"github.com/k-tahiro/tahiremogon/model"
)

const (
	height = 224
	width  = 224
)

var (
	mean = [3]float32{0.485, 0.456, 0.406}
	std  = [3]float32{0.229, 0.224, 0.225}
)

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
	if _, err := exec.Command("sh", "-c", cmd).Output(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	cmd = "sudo ../bin/camera.sh"
	filename, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	input, err := readImage(*(*string)(unsafe.Pointer(&filename)))
	predict(cc.PredictionModel, input)

	var response model.Response
	response.Success = true
	return cc.JSON(http.StatusOK, response)
}

func readImage(filename string) (tensor.Tensor, error) {
	// Read input image
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	img = resize.Resize(width, height, img, resize.Bilinear)

	input := tensor.New(tensor.WithShape(1, 3, height, width), tensor.Of(tensor.Float32))
	err = imageToBCHW(img, input)
	if err != nil {
		return nil, err
	}
	err = normalize(input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func imageToBCHW(img image.Image, dst tensor.Tensor) error {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, g, b, a := img.At(x, y).RGBA()
			if a != 65535 {
				return errors.New("transparency not handled")
			}
			err := dst.SetAt(float32(uint8(r/0x100)), 0, 0, y, x)
			if err != nil {
				return err
			}
			err = dst.SetAt(float32(uint8(g/0x100)), 0, 1, y, x)
			if err != nil {
				return err
			}
			err = dst.SetAt(float32(uint8(b/0x100)), 0, 2, y, x)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func normalize(input tensor.Tensor) (err error) {
	for channel := 0; channel < 3; channel++ {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				z, err := input.At(0, channel, x, y)
				if err != nil {
					return err
				}
				zn := z.(float32) / 255
				err = input.SetAt((zn-mean[channel])/std[channel], 0, channel, x, y)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func predict(predictionModel *myMw.PredictionModel, input tensor.Tensor) (int, error) {
	backend := predictionModel.Graph
	m := predictionModel.Model

	m.SetInput(0, input)
	err := backend.Run()
	if err != nil {
		return -1, err
	}
	output, err := m.GetOutputTensors()
	if err != nil {
		return -1, err
	}

	// Find maximum value of prediction results
	max := float32(-9999)
	maxi := -1
	for i, v := range output[0].Data().([]float32) {
		fmt.Println(i, v)
		if v > max {
			max = v
			maxi = i
		}
	}

	return maxi, nil
}
