package http

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

type Http struct {
	e    *echo.Echo
	port uint16
}

func NewHttp(debug bool, port uint16) *Http {
	e := echo.New()
	e.HideBanner = true

	e.GET("/metrics", echoprometheus.NewHandler())
	if debug {
		e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))
	}

	return &Http{
		e,
		port,
	}
}

func (h *Http) Serve() error {
	log.Printf("Starting http server on port %d\n", h.port)
	return h.e.Start(":" + strconv.Itoa(int(h.port)))
}

func (h *Http) Stop(ctx context.Context) error {
	return h.e.Shutdown(ctx)
}
