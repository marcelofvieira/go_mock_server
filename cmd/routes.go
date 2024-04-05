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
	router := mux.NewRouter()

	// ----------------------------------------------------------------------------------------------------------------
	// Mock Configuration Handler
	// ----------------------------------------------------------------------------------------------------------------
	repository := memorykvs.NewMemKVS()
	mockRequestPreProcessor := requestpreprocessor.NewService()
	mockConfiguration := mockconfiguration.NewService(repository, mockRequestPreProcessor)

	mockConfigHandler := mockconfighandler.NewHTTPHandler(mockConfiguration)

	router.HandleFunc("/api/v1/configurations/{id}", mockConfigHandler.GetMockConfiguration).Methods("GET")
	router.HandleFunc("/api/v1/configurations", mockConfigHandler.AddMockConfiguration).Methods("POST")
	router.HandleFunc("/api/v1/configurations/{id}", mockConfigHandler.UpdateMockConfiguration).Methods("PUT")
	router.HandleFunc("/api/v1/configurations/{id}", mockConfigHandler.DeleteMockConfiguration).Methods("DELETE")

	// ----------------------------------------------------------------------------------------------------------------
	// Dynamic Handler
	// ----------------------------------------------------------------------------------------------------------------
	filterHandler := requestfilter.NewService()
	variableProcessor := variableprocessor.NewService()
	responseProcessor := responseprocessor.NewService()
	dynamicHandlerProcessor := dynamichandlerprocessor.NewService(repository, filterHandler, variableProcessor, responseProcessor)

	dynamicHandler := dynamichandler.NewHTTPHandler(dynamicHandlerProcessor)

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
