package main

import (
	"fmt"

	"github.com/url-shortener/config"
)

func main() {
	config.Init()

	fmt.Println(config.Log())
	return
}
