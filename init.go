package dockerlib

import "go.uber.org/zap"

var logger *zap.SugaredLogger
var instance *Controller

func init() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger = log.Named("docker").Sugar()

	instance, err = NewDockerController()
	if err != nil {
		panic(err)
	}
}
