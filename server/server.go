package server

import (
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/thehivecorporation/raccoon/constants"
	"github.com/thehivecorporation/raccoon/parser"
)

type response struct {
	Status string `json:"status"`
}

func REST(c *cli.Context) {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", func(c echo.Context) error {

		byt, err := ioutil.ReadAll(c.Request().Body())
		if err != nil {
			return err
		}

		err = parser.ParseRequest(byt)
		if err != nil {
			log.Error(err)
			return err
		}

		s := response{
			Status: "ok",
		}

		return c.JSON(http.StatusOK, s)
	})

	log.WithFields(log.Fields{
		constants.PORT_FLAG_NAME: c.String("port"),
	}).Info("Starting Raccoon server...")

	e.Run(standard.New(":" + c.String("port")))
}
