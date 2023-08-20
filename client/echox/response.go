package echox

import (
	"github.com/labstack/echo/v4"
)

// Response in order to unify the returned response structure
type Response struct {
	Code    int    `json:"-"`
	Pretty  bool   `json:"-"`
	Status  int    `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

// JSON sends a JSON response with status code.
func (a Response) JSON(ctx echo.Context) error {
	a.Status = 1
	if len(a.Message) > 0 {
		a.Status = -1
	}
	if a.Pretty {
		return ctx.JSONPretty(a.Code, a, "\t")
	}
	return ctx.JSON(a.Code, a)
}
