package server

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"github.com/thehivecorporation/raccoon"
	"github.com/thehivecorporation/raccoon/parser"
)

//REST is the server that is launched when a user selects the "server" option
//in the CLI
func REST(c *cli.Context) {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/", func(c echo.Context) error {
		jobParser := parser.Job{Dispatcher: new(raccoon.SimpleDispatcher)}

		var jobRequest raccoon.JobRequest
		if err := jobParser.Build(c.Request().Body(), &jobRequest); err != nil {
			return err
		}

		taskList, err := jobParser.ParseTaskList(jobRequest.TaskList)
		if err != nil {
			return err
		}

		jobs := jobParser.BuildJobList(jobRequest.Infrastructure, taskList)

		//Send jobs to dispatcher
		go jobParser.Dispatcher.Dispatch(*jobs)

		rsp := struct {
			Status string
		}{
			Status: "ok",
		}

		return c.JSON(http.StatusOK, rsp)
	})

	log.WithFields(log.Fields{
		"port": c.String("port"),
	}).Info("Starting Raccoon server...")

	e.Run(standard.New(":" + c.String("port")))
}
