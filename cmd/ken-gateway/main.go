package main

import (
	"ken-gateway/internal/ken-gateway/router"
)

const (
	host = "0.0.0.0"
	port = "9999"
)

func main() {
	r := router.InitRouter()

	r.Run(host + ":" + port)
}
