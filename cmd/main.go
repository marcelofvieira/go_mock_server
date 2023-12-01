package main

import (
	_ "encoding/json"
	_ "mock_server_mux/pkg/apperrors"
	"mock_server_mux/pkg/logger"
)

func main() {
	err := SetRoutes()

	if err != nil {
		logger.Error("Error: ", err)
		panic(err)
	}
}
