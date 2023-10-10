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

	mockConfigRepository := mockconfigurationrepo.NewMemKVS()
	configMockService := mockconfiguration.NewService(mockConfigRepository)
	mockConfigHandler := mockconfighandler.NewHTTPHandler(configMockService)

	router.HandleFunc("/mock-config/{id}", mockConfigHandler.GetMockConfiguration).Methods("GET")
	router.HandleFunc("/mock-config", mockConfigHandler.AddMockConfiguration).Methods("POST")
	router.HandleFunc("/mock-config", mockConfigHandler.UpdateMockConfiguration).Methods("PUT")
	router.HandleFunc("/mock-config/{id}", mockConfigHandler.DeleteMockConfiguration).Methods("DELETE")

	filterHandlerService := requestfilter.NewService()
	processorService := dynamichandlerprocessor.NewService(mockConfigRepository, filterHandlerService)
	dynamicHandler := dynamichandler.NewHTTPHandler(processorService)

	router.NotFoundHandler = http.HandlerFunc(dynamicHandler.ProcessDynamicHandler)

	http.Handle("/", router)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}

	return nil

}

/*
func getMockByIdHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GET BY Id")
}

func getMockHandler(w http.ResponseWriter, r *http.Request) {
	FormatResponse(w, r, http.StatusOK, customRoutes)
}

func genericHandler(w http.ResponseWriter, r *http.Request, configMock domain.MockConfiguration) {

	FormatResponse(w, r, configMock.Response.StatusCode, configMock.Response.Body)

}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	route := r.Method + " " + r.URL.Path

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyString := string(bodyBytes)

	fmt.Println(bodyString)

	for _, configMock := range customRoutes {

		pattern := configMock.Request.Method + " " + configMock.Request.URL

		validRegex := regexp.MustCompile(pattern)

		if validRegex.MatchString(route) {
			genericHandler(w, r, configMock)
			return
		}
	}

	FormatResponse(w, r, http.StatusBadRequest, "{'status': 'not found'}")

	//http.NotFoundHandler().ServeHTTP(w, r)

}
*/
