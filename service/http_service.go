package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"disfactory/imgur-backup/utils"
)

type HttpService struct {
	address   string
	backupDir string
}

func NewHttpService(address string, backupDir string) *HttpService {
	return &HttpService{
		address:   address,
		backupDir: backupDir,
	}
}

func (s HttpService) Start() {
	e := echo.New()
	echoLog := utils.EchoWrapperLogger()
	e.Logger = echoLog

	e.Use(middleware.Logger())

	imgurBackGroup := e.Group("/imgur")
	imgurBackGroup.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: s.backupDir,
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(s.address))
}
