package dockerlib

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func init() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger = log.Named("docker").Sugar()
}
