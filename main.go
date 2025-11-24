package main

import (
	"github.com/lemonc7/project-manager/application"
)

func main() {
	app, err := application.InitApp("./config.yaml")
	if err != nil {
		panic(err)
	}

	app.Start()
}
