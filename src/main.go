package main

import (
	"github.com/url-shortener/config"
	"github.com/url-shortener/logger"
)

func main() {
	config.Init()

	logger.Init(config.Log())
	defer logger.Flush()

	logger.Debug("Hello debug")
	logger.Info("Hello info")

	return
}
