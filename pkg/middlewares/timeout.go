package middlewares

import (
	"context"
	"github.com/labstack/echo/v4"
	"time"
)

func Timeout(timeout time.Duration) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			go func() {
				<-ctx.Done()
				cancel()
			}()
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
