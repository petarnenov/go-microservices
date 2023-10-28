package main

import (
	"logger/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var logEntry JSONPayload

	err := app.readJSON(w, r, &logEntry)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//insert data
	event := data.LogEntry{
		Name: logEntry.Name,
		Data: logEntry.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "log entry inserted",
	}

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
}
