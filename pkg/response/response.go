package response

import (
	"encoding/json"
	"log"
	"mock_server_mux/api/presenter"
	"mock_server_mux/pkg/apperrors"
	"net/http"
	"time"
)

func ProcessMockResponse(resp http.ResponseWriter, req *http.Request, httpStatus int, headers map[string]string, body interface{}, delay int) error {

	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}

	flusher, ok := resp.(http.Flusher)
	if !ok {
		http.NotFound(resp, req)
		return nil
	}

	for key, value := range headers {
		resp.Header().Set(key, value)
	}

	resp.WriteHeader(httpStatus)

	if body != nil {
		jsonResp, err := json.Marshal(body)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}

		if _, err := resp.Write(jsonResp); err != nil {
			return err
		}
	} else {
		resp.Write(nil)
	}

	flusher.Flush()

	return nil
}

func Success(w http.ResponseWriter, r *http.Request, statusCode int, body interface{}) error {

	jsonResp, err := json.Marshal(body)
	if err != nil {
		return err
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, r)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonResp); err != nil {
		return err
	}

	flusher.Flush()

	return nil
}

func Error(w http.ResponseWriter, r *http.Request, statusCode int, err error) error {

	responseError := presenter.ResponseError{
		Code:    apperrors.Code(err),
		Message: err.Error(),
	}

	jsonResp, err := json.Marshal(responseError)
	if err != nil {
		return err
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.NotFound(w, r)
		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonResp); err != nil {
		return err
	}

	flusher.Flush()

	return nil
}
