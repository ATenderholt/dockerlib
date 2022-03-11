package dockerlib

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func init() {
	//config := zap.NewProductionConfig()
	//config.Encoding = "console"
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	logger = log.Named("docker").Sugar()
}

func SetLogger(newLogger *zap.Logger) {
	logger = newLogger.Named("docker").Sugar()
}
