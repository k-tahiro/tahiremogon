package middleware

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
	"gorgonia.org/tensor"
)

type PredictionModel struct {
	Graph *gorgonnx.Graph
	Model *onnx.Model
}

func loadPredictionModel(model string) (*PredictionModel, error) {
	backend := gorgonnx.NewGraph()
	m := onnx.NewModel(backend)

	// Read model binary
	b, err := ioutil.ReadFile(model)
	if err != nil {
		return nil, err
	}

	// Decode it into the model
	err = m.UnmarshalBinary(b)
	if err != nil {
		return nil, err
	}

	predictionModel := &PredictionModel{
		Graph: backend,
		Model: m,
	}
	return predictionModel, nil
}

func PredictionModelMiddleware(modelFile string) echo.MiddlewareFunc {
	predictionModel, err := loadPredictionModel(modelFile)
	if err != nil {
		log.Fatal(err)
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cctx, ok := c.(*CustomContext)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "カスタムコンテキストが取得できません")
			}

			cctx.PredictionModel = predictionModel
			return next(cctx)
		}
	}
}

func (predictionModel *PredictionModel) Predict(input tensor.Tensor) (int, error) {
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
