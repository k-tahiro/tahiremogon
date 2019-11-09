package main

import (
	"log"
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
	model, err := myMw.LoadPredictionModel(os.Getenv("ONNX_MODEL_FILE"))
	if err != nil {
		log.Fatal(err)
	}
	e.Use(myMw.PredictionModelMiddleware(model))

	// Routes
	commands := e.Group("/commands")
	{
		commands.GET("/", handler.Commands)
		commands.POST("/transmit/:id", handler.Transmit)
		commands.POST("/receive", handler.Receive)
	}

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
