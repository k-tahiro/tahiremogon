package middleware

import (
	"github.com/gocraft/dbr/v2"
	"github.com/labstack/echo"
)

// カスタムコンテキスト
type CustomContext struct {
	echo.Context
	Connection      *dbr.Connection
	PredictionModel *PredictionModel
}

// カスタムコンテキストを定義するMiddleware
func MyCustomContextMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// カスタムコンテキストを初期化して次へ
			cctx := &CustomContext{
				Context: c,
			}
			return next(cctx)
		}
	}
}
