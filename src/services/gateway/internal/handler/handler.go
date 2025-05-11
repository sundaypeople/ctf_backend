package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/LainInTheWired/ctf_backend/gateway/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

// 依存関係用の構造体
type GatewayHandler struct {
	serv service.GatewayService
}

func NewHanderGateway(service service.GatewayService) *GatewayHandler {
	return &GatewayHandler{
		serv: service,
	}
}

func (h *GatewayHandler) Authz(c echo.Context) error {
	id, err := c.Cookie("session")
	if err != nil {
		wrappedErr := errors.Wrap(err, "request bind error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": fmt.Sprintf("not login")})
	}
	u, err := h.serv.GetUserID(id.Value)
	if err != nil {
		wrappedErr := errors.Wrap(err, "request bind error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": fmt.Sprintf("not login")})
	}
	rs, err := h.serv.GetRoles(u)
	if err != nil {
		wrappedErr := errors.Wrap(err, "request bind error")
		log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": fmt.Sprintf("error: %v", wrappedErr)})
	}
	log.Print(rs)
	xForwardedUri := c.Request().Header.Get("X-Frontend-Path")
	log.Printf("%#v", c.Request().Header)
	if xForwardedUri != "" {
		segments := strings.Split(xForwardedUri, "/")
		fmt.Println(segments)
		if segments[1] == "admin" {
			fmt.Print("fewa")
			for _, r := range rs {
				if r.ID != 1 {
					wrappedErr := errors.Wrap(err, "request bind error")
					log.Errorf("\n%+v\n", wrappedErr) // スタックトレース付きでログに出力
					return c.JSON(http.StatusForbidden, map[string]string{"message": "success bind user roles"})
				}
			}
		}
	}

	c.Response().Header().Set("X-User-ID", strconv.Itoa(u))
	fmt.Println("X-User-ID", strconv.Itoa(u))

	return c.JSON(http.StatusCreated, map[string]string{"message": "success bind user roles"})
}
