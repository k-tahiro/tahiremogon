package middleware

import (
	"fmt"
	"net/http"

	"github.com/gocraft/dbr"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3" // まずは既存のSQLiteのデータベースへ接続する
	"golang.org/x/crypto/ssh"
)

// カスタムコンテキスト
type CustomContext struct {
	echo.Context
	Connection *dbr.Connection
	Client     *ssh.Client
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

// SSHクライアントを設定するMiddleware
func SSHClientMiddleware(address string, user string, password string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cctx, ok := c.(*CustomContext)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "カスタムコンテキストが取得できません")
			}

			var hostKey ssh.PublicKey
			config := &ssh.ClientConfig{
				User: user,
				Auth: []ssh.AuthMethod{
					ssh.Password(password),
				},
				HostKeyCallback: ssh.FixedHostKey(hostKey),
			}

			hostport := fmt.Sprintf("%s:%d", address, 22)
			client, err := ssh.Dial("tcp", hostport, config)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "SSHが接続できません")
			}

			// リモコンへのSSH接続をコンテキストに設定して次へ
			cctx.Client = client
			return next(cctx)
		}
	}
}
