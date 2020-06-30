package main

import (
	"github.com/vladikan/url-shortener/config"
	"github.com/vladikan/url-shortener/logger"
	"github.com/vladikan/url-shortener/service"
)

func main() {
	config.Init("")

	logger.Init(config.Log())
	defer logger.Flush()

	service.Start(config.Host())
}
