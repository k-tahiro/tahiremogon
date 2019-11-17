package middleware

import (
	"net/http"

	"github.com/gocraft/dbr/v2"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3" // まずは既存のSQLiteのデータベースへ接続する
)

// SQLiteを設定するMiddleware
func SQLiteMiddleware(datasource string) echo.MiddlewareFunc {
	conn, err := dbr.Open("sqlite3", datasource, nil)
	if err != nil {
		panic(err)
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cctx, ok := c.(*CustomContext)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "カスタムコンテキストが取得できません")
			}

			cctx.Session = conn.NewSession(nil)
			return next(cctx)
		}
	}
}
