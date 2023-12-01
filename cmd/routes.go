package main

import (
	"github.com/gorilla/mux"
	"mock_server_mux/api/handlers/dynamichandler"
	"mock_server_mux/api/handlers/mockconfighandler"
	"mock_server_mux/internal/core/service/dynamichandlerprocessor"
	"mock_server_mux/internal/core/service/mockconfiguration"
	"mock_server_mux/internal/core/service/requestfilter"
	"mock_server_mux/internal/core/service/requestpreprocessor"
	"mock_server_mux/internal/core/service/responseprocessor"
	"mock_server_mux/internal/core/service/variableprocessor"
	"mock_server_mux/internal/repository/memorykvs"
	"net/http"
)

func SetRoutes() error {
	//TODO: Define best place for router
	router := mux.NewRouter()

	// ----------------------------------------------------------------------------------------------------------------
	// Mock Configuration Handler
	// ----------------------------------------------------------------------------------------------------------------
	mockConfigRepository := memorykvs.NewMemKVS()
	mockRequestPreProcessorService := requestpreprocessor.NewService()
	configMockService := mockconfiguration.NewService(mockConfigRepository, mockRequestPreProcessorService)

	mockConfigHandler := mockconfighandler.NewHTTPHandler(configMockService)

	router.HandleFunc("/mock-config/{id}", mockConfigHandler.GetMockConfiguration).Methods("GET")
	router.HandleFunc("/mock-config", mockConfigHandler.AddMockConfiguration).Methods("POST")
	router.HandleFunc("/mock-config/{id}", mockConfigHandler.UpdateMockConfiguration).Methods("PUT")
	router.HandleFunc("/mock-config/{id}", mockConfigHandler.DeleteMockConfiguration).Methods("DELETE")

	// ----------------------------------------------------------------------------------------------------------------
	// Dynamic Handler
	// ----------------------------------------------------------------------------------------------------------------
	filterHandlerService := requestfilter.NewService()
	responseProcessorService := responseprocessor.NewService()
	mockProcessorService := variableprocessor.NewService()
	processorService := dynamichandlerprocessor.NewService(mockConfigRepository, filterHandlerService, mockProcessorService, responseProcessorService)

	dynamicHandler := dynamichandler.NewHTTPHandler(processorService)

	router.NotFoundHandler = http.HandlerFunc(dynamicHandler.ProcessDynamicHandler)

	// ----------------------------------------------------------------------------------------------------------------
	// Init server
	// ----------------------------------------------------------------------------------------------------------------
	http.Handle("/", router)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}

	return nil
}
