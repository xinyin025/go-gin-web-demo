package main

import (
	"go-web-demo/config"
	"go-web-demo/routes"
)

func main() {
	config.ConnectDatabase()

	r := routes.SetupRouter()

	err := r.Run(":8082")
	if err != nil {
		return
	}

}
