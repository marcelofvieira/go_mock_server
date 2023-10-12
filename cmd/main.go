package main

import (
	_ "encoding/json"
	"github.com/gorilla/mux"
	"mock_server_mux/api/handlers/dynamichandler"
	"mock_server_mux/api/handlers/mockconfighandler"
	"mock_server_mux/internal/core/service/dynamichandlerprocessor"
	"mock_server_mux/internal/core/service/mockconfiguration"
	"mock_server_mux/internal/core/service/requestfilter"
	"mock_server_mux/internal/repository/mockconfigurationrepo"
	_ "mock_server_mux/pkg/apperrors"
	"mock_server_mux/pkg/logger"
	"net/http"
)

func main() {
	err := Run()

	if err != nil {
		logger.Error("Error: ", err)
		panic(err)
	}
}

func Run() error {

	router := mux.NewRouter()

	// ----------------------------------------------------------------------------------------------------------------
	// Mock Configuration Handler
	// ----------------------------------------------------------------------------------------------------------------
	mockConfigRepository := mockconfigurationrepo.NewMemKVS()
	configMockService := mockconfiguration.NewService(mockConfigRepository)
	mockConfigHandler := mockconfighandler.NewHTTPHandler(configMockService)

	router.HandleFunc("/mock-config/{id}", mockConfigHandler.GetMockConfiguration).Methods("GET")
	router.HandleFunc("/mock-config", mockConfigHandler.AddMockConfiguration).Methods("POST")
	router.HandleFunc("/mock-config", mockConfigHandler.UpdateMockConfiguration).Methods("PUT")
	router.HandleFunc("/mock-config/{id}", mockConfigHandler.DeleteMockConfiguration).Methods("DELETE")

	// ----------------------------------------------------------------------------------------------------------------
	// Dynamic Handler
	// ----------------------------------------------------------------------------------------------------------------
	filterHandlerService := requestfilter.NewService()
	processorService := dynamichandlerprocessor.NewService(mockConfigRepository, filterHandlerService)
	dynamicHandler := dynamichandler.NewHTTPHandler(processorService)

	router.NotFoundHandler = http.HandlerFunc(dynamicHandler.ProcessDynamicHandler)

	// ----------------------------------------------------------------------------------------------------------------
	// Dynamic Handler
	// ----------------------------------------------------------------------------------------------------------------
	http.Handle("/", router)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}

	return nil
}
