package pkg

import (
	"encoding/json"
	"fmt"
	pkgErr "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vvinokurshin/AvitoInternship/pkg/errors"
	"net/http"
	"os"
	"time"
)

type ResponseWriterCode struct {
	http.ResponseWriter
	StatusCode int
}

func NewResponseWriterCode(w http.ResponseWriter) *ResponseWriterCode {
	return &ResponseWriterCode{w, http.StatusOK}
}

func (rw *ResponseWriterCode) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	causeErr := pkgErr.Cause(err)
	code := errors.HttpCode(causeErr)
	customErr := errors.New(code, causeErr)
	logLevel := errors.LogLevel(causeErr)

	globalLogger, ok := r.Context().Value(ContextHandlerLog).(*Logger)
	if !ok {
		log.Error("failed to get logger for handler", r.URL.Path)
		log.Error(err)
	} else {
		globalLogger.Log(logLevel, err)
	}

	SendJSON(w, r, code, customErr)
}

func SendJSON(w http.ResponseWriter, r *http.Request, status int, dataStruct any) {
	dataJSON, err := json.Marshal(dataStruct)
	if err != nil {
		HandleError(w, r, fmt.Errorf("failed to marshal : %w", err))
		return
	}

	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(status)

	_, err = w.Write(dataJSON)
	if err != nil {
		HandleError(w, r, fmt.Errorf("failed to send : %w", err))
		return
	}
}

func SendFile(w http.ResponseWriter, r *http.Request, fileName string) {
	file, err := os.Open(HistoryFolderName + fileName)
	if err != nil {
		HandleError(w, r, pkgErr.Wrap(err, "file not found"))
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	http.ServeContent(w, r, fileName, time.Now(), file)
}
