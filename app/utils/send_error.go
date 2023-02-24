package utils

import (
	"github.com/labstack/echo/v4"
)

func SendError(c echo.Context, status int, err string) error {
	return c.JSON(status, err)
}
