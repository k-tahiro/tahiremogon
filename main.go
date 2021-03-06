package main

import (
	"os"

	"github.com/labstack/echo"
	echoMw "github.com/labstack/echo/middleware"

	"github.com/k-tahiro/tahiremogon/handler"
	myMw "github.com/k-tahiro/tahiremogon/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(echoMw.Logger())
	e.Use(echoMw.Recover())

	// カスタムコンテキスト用Middlewareを適用
	e.Use(myMw.MyCustomContextMiddleware())

	// DB用Middlewareを適用
	e.Use(myMw.SQLiteMiddleware(os.Getenv("DB_FILE")))

	// エアコン状態判定用モデルMiddlewareを適用
	e.Use(myMw.PredictionModelMiddleware(os.Getenv("ONNX_MODEL_FILE")))

	// Routes
	codes := e.Group("/codes")
	{
		codes.GET("/", handler.ReadCodes)
		codes.POST("/:key", handler.CreateCode)
		codes.DELETE("/:key", handler.DeleteCode)
		codes.POST("/transmit/:key", handler.TransmitCode)
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
