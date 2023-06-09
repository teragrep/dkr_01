package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/teragrep/dkr_02/sdk"
	"log"
	"os"
)

var logLevels = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"info":  logrus.InfoLevel,
	"warn":  logrus.WarnLevel,
	"error": logrus.ErrorLevel,
}

func main() {
	fmt.Fprintln(os.Stdout, "RELP Docker Plugin start")
	selectedLogLevel := os.Getenv("LOG_LEVEL")
	if selectedLogLevel == "" {
		// default to INFO log level if no LOG_LEVEL environment variable set
		selectedLogLevel = "info"
	}

	if level, exists := logLevels[selectedLogLevel]; exists {
		logrus.SetLevel(level)
	} else {
		log.Println("Invalid logging level: " + selectedLogLevel)
		os.Exit(1)
	}

	relpPluginHandler := sdk.NewHandler(`{"Implements": ["LogDriver"]}`)
	handlers(&relpPluginHandler, newDriver())
	if err := relpPluginHandler.ServeUnix("dkr_01", 0); err != nil {
		panic(err)
	}
}
