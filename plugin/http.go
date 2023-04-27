package main

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/daemon/logger"
	"github.com/teragrep/dkr_02/sdk"
	"io"
	"net/http"
	"os"
)

type StartLoggingRequest struct {
	File string
	Info logger.Info
}

type StopLoggingRequest struct {
	File string
}

type CapabilitiesResponse struct {
	Err string
	Cap logger.Capability
}

type ReadLogsRequest struct {
	Info   logger.Info
	Config logger.ReadConfig
}

type response struct {
	Err string
}

func handlers(h *sdk.Handler, d *Driver) {
	h.HandleFunc("/LogDriver.StartLogging", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req StartLoggingRequest
		fmt.Fprintf(os.Stdout, "Start logging request for container %s", body)
		json.Unmarshal(body, &req)
		err := (*d).StartLogging(req.File, req.Info)
		respond(err, w)
	})

	h.HandleFunc("/LogDriver.StopLogging", func(w http.ResponseWriter, r *http.Request) {
		var req StopLoggingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprintf(os.Stdout, "Stop logging request was called for container")
		err := d.StopLogging(req.File)
		respond(err, w)
	})

	h.HandleFunc("/LogDriver.Capabilities", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(&CapabilitiesResponse{
			Cap: logger.Capability{ReadLogs: false},
		})
	})

	h.HandleFunc("/LogDriver.ReadLogs", func(w http.ResponseWriter, r *http.Request) {
		var req ReadLogsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		fmt.Fprintln(os.Stdout, "ReadLogs called for container ", req.Info.ContainerID)
		http.Error(w, "not impl", http.StatusNotImplemented)
	})
}

func respond(err error, w http.ResponseWriter) error {
	var res response
	if err != nil {
		res.Err = err.Error()
	}
	return json.NewEncoder(w).Encode(&res)
}
