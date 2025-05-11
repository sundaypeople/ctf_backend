package main

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/xerrors"

	"github.com/LainInTheWired/ctf_backend/gateway/handler"
	"github.com/LainInTheWired/ctf_backend/gateway/repository"
	"github.com/LainInTheWired/ctf_backend/gateway/service"
	_ "github.com/go-sql-driver/mysql" // 空のインポートを追加

	"github.com/gorilla/sessions"

	// 正しいモジュールパス
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoLog "github.com/labstack/gommon/log" // エイリアスを付ける
)

type RedisConfig struct {
	IP       string `env:"REDIS_IP"`
	Port     string `env:"REDIS_PORT"`
	Password string `env:"REDIS_PASSWORD"`
}

func main() {

	reddb, err := repository.NewRedis()
	if err != nil {
		xerrors.Errorf("redis connetciono error: %w", err.Error())
	}
	defer reddb.Close()

	e := echo.New()
	// セッション
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	// ログの表示方式の切り替え
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}
	// ログフォーマットの設定
	if env == "production" {
		// 本番環境ではJSON形式のログを出力
		echoLog.SetHeader(`{"time":"${time_rfc3339}","level":"${level}","file":"${short_file}","line":"${line}","message":"${message}"}`)
		echoLog.SetLevel(echoLog.ERROR) // 必要に応じてログレベルを調整
	} else {
		// 開発環境ではテキスト形式のログを出力
		echoLog.SetHeader("${time_rfc3339} [${level}] ${short_file}:${line} ${message}")
		echoLog.SetLevel(echoLog.DEBUG) // 開発中は詳細なログを出力
	}

	// アクセスログを出力
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	rr := repository.NewRedisClient(reddb, context.Background())
	s := service.NewGatewayService(rr)
	h := handler.NewHanderGateway(s)

	// 認可エンドポイントの定義
	// e.GET("/auth", func(c echo.Context) error {
	// })
	e.GET("/auth", h.Authz)

	e.Start(":8000")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
