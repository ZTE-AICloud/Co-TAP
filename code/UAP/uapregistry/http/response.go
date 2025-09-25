package http

import (
	"fmt"
	"net/http"
	"strconv"
	"uapregistry/logger"
)

func Response200(w http.ResponseWriter, jsonbytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if _, err := w.Write(jsonbytes); err != nil {
		logger.GetLogger().Errorf("Failed to write 200 on the connection:%v", err)
	}
}

func Response200WithContentType(w http.ResponseWriter, contentType string, jsonbytes []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(200)
	if _, err := w.Write(jsonbytes); err != nil {
		logger.GetLogger().Errorf("Failed to write 200 on the connection:%v", err)
	}
}

func Response200WithCustomHeader(w http.ResponseWriter, jsonbytes []byte, customHeader map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	for k, v := range customHeader {
		w.Header().Set(k, v)
	}
	w.WriteHeader(200)
	if _, err := w.Write(jsonbytes); err != nil {
		logger.GetLogger().Errorf("Failed to write 200custom on the connection:%v", err)
	}
}

func Response201(w http.ResponseWriter, jsonbytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	if _, err := w.Write(jsonbytes); err != nil {
		logger.GetLogger().Errorf("Failed to write 201 on the connection:%v", err)
	}
}

func Response201WithContentType(w http.ResponseWriter, contentType string, jsonbytes []byte) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(201)
	if _, err := w.Write(jsonbytes); err != nil {
		logger.GetLogger().Errorf("Failed to write 201 on the connection:%v", err)
	}
}

func Response204(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(204)
}

func Response400(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)
	if _, err := fmt.Fprintln(w, `{"code":"400","message":"`+message+`"}`); err != nil {
		logger.GetLogger().Errorf("Failed to write 400 on the connection:%v", err)
	}
}

func Response404(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	if _, err := fmt.Fprintln(w, `{"code":"404","message":"`+message+`"}`); err != nil {
		logger.GetLogger().Errorf("Failed to write 404 on the connection:%v", err)
	}
}

func Response404WithCustomHeader(w http.ResponseWriter, message string, customHeader map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	for k, v := range customHeader {
		w.Header().Set(k, v)
	}
	w.WriteHeader(404)
	if _, err := fmt.Fprintln(w, `{"code":"404","message":"`+message+`"}`); err != nil {
		logger.GetLogger().Errorf("Failed to write 404custom on the connection:%v", err)
	}
}

func Response408(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(408)
	fmt.Fprintln(w, `{"code":"408","message":"request timeout"}`)
}

func Response409(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)
	if _, err := fmt.Fprintln(w, `{"code":"409","message":"`+message+`"}`); err != nil {
		logger.GetLogger().Errorf("Failed to write 409 on the connection:%v", err)
	}
}

func Response415(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(415)
	if _, err := fmt.Fprintln(w, `{"code":"415","message":"`+message+`"}`); err != nil {
		logger.GetLogger().Errorf("Failed to write 415 on the connection:%v", err)
	}
}

func Response422(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(422)
	if _, err := fmt.Fprintln(w, `{"code":"422","message":"`+message+`"}`); err != nil {
		logger.GetLogger().Errorf("Failed to write 422 on the connection:%v", err)
	}
}

func Response500(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	if _, err := fmt.Fprintln(w, `{"code":"500","message":"`+message+`"}`); err != nil {
		logger.GetLogger().Errorf("Failed to write 500 on the connection:%v", err)
	}
}

func ResponseWithStatusCode(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := fmt.Fprintln(w, `{"code":"`+strconv.Itoa(statusCode)+`","message":"`+message+`"}`); err != nil {
		logger.GetLogger().Errorf("Failed to write %s on the connection:%v", statusCode, err)
	}
}
