package gateway

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) CheckServer(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, nil)
}
