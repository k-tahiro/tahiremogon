package middleware

import (
	"net/http"
	"github.com/labstack/echo"
	"github.com/gocraft/dbr"
	_ "github.com/mattn/go-sqlite3" // まずは既存のSQLiteのデータベースへ接続する
	_ "github.com/go-sql-driver/mysql" // 本当はmain.goに置いほうが良いらしい
)

// カスタムコンテキスト
type CustomContext struct {
    echo.Context
    Connection            *dbr.Connection
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

/*
// MySQLを設定するMiddleware
func MySQLMiddleware(datasource string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            cctx, ok := c.(*CustomContext)
            if !ok {
                return echo.NewHTTPError(http.StatusInternalServerError, "カスタムコンテキストが取得できません")
            }

            db, err := sql.Open("mysql", datasource)
            if err != nil {
                return echo.NewHTTPError(http.StatusInternalServerError, "DBが取得できません")              
            }
            defer db.Close()

            // DBをコンテキストに設定して次へ
            cctx.DB = db
            return next(cctx)
        }
    }
}
*/

// SQLiteを設定するMiddleware
func SQLiteMiddleware(datasource string) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            cctx, ok := c.(*CustomContext)
            if !ok {
                return echo.NewHTTPError(http.StatusInternalServerError, "カスタムコンテキストが取得できません")
            }

            conn, err := dbr.Open("sqlite3", datasource, nil)
            if err != nil {
                return echo.NewHTTPError(http.StatusInternalServerError, "DBが取得できません")              
            }

            // DBへのConnectionをコンテキストに設定して次へ
            cctx.Connection = conn
            return next(cctx)
        }
    }
}
