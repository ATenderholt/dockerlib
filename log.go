package dockerlib

import "go.uber.org/zap"

func SetLogger(newLogger *zap.Logger) {
	logger = newLogger.Named("docker").Sugar()
}
