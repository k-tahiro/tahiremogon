package middleware

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/owulveryck/onnx-go"
	"github.com/owulveryck/onnx-go/backend/x/gorgonnx"
)

type PredictionModel struct {
	Graph *gorgonnx.Graph
	Model *onnx.Model
}

func LoadPredictionModel(model string) (*PredictionModel, error) {
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

	var predictionModel *PredictionModel
	predictionModel.Graph = backend
	predictionModel.Model = m
	return predictionModel, nil
}

func PredictionModelMiddleware(predictionModel *PredictionModel) echo.MiddlewareFunc {
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
