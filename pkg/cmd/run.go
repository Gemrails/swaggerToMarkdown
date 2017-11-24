package cmd

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/akkuman/parseConfig"
)

//SwaggerAction SwaggerAction
type SwaggerAction struct{}

//CreateSWManager CreateManager
func CreateSWManager() *SwaggerAction {
	return &SwaggerAction{}
}

//ShowConf ShowConf
func (s *SwaggerAction) ShowConf(sw string) error {
	if err := checkConf(sw); err != nil {
		logrus.Errorf("config check error, %v", err)
		return err
	}
	config := getConf(sw)
	logrus.Infof("file is %v", config.Get("/license"))
	return nil
}

func checkConf(confPath string) error {
	if _, err := os.Stat(confPath); os.IsNotExist(err) {
		return fmt.Errorf("config.json is not exist")
	}
	return nil
}

func getConf(confPath string) parseConfig.Config {
	return parseConfig.New(confPath)
}
