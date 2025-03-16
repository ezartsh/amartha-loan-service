package main

import (
	"fmt"
	"loan-service/config"
	"loan-service/pkg/logger"
	"loan-service/pkg/mail"
	"loan-service/routes"
	"net/http"
	"os"
	"time"
)

func main() {
	startTime := time.Now()
	logger.AppLog.Debug("Service Init")

	configResults, errConfigBags, err := config.Init()
	logger.Register()

	if err != nil {
		logger.AppLog.Error(err, "[Env]")
		for _, err := range errConfigBags {
			logger.AppLog.DebugWithVariables(logger.LogVariable{
				"name":        err.VariableName,
				"error":       err.Message,
				"description": err.Description,
			}, "[Env]")
		}
		os.Exit(1)
	} else {
		logger.AppLog.DebugWithVariables(configResults, "[Env] Setup Configuration Value")
	}

	mailHandler := mail.NewMail(mail.NewMailTrap())

	logger.AppLog.Info("Initiate route.")
	routeHandler := routes.WebInit(mailHandler)

	logger.AppLog.Infof("Running service on port %d.", config.Env.HTTPPort)

	logger.AppLog.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%d", config.Env.HTTPPort),
		routeHandler,
	))

	logger.AppLog.Info("Init finished in ", time.Since(startTime))
	logger.AppLog.Info("Service Ready")

}
