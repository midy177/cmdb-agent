package handler

import (
	"github.com/labstack/echo/v4"
)

type StreamWriter struct {
	W *echo.Response
}

func (s *StreamWriter) Write(p []byte) (n int, err error) {
	defer s.W.Flush()
	return s.W.Write(p)
}
