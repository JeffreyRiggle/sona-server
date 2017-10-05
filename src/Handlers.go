package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var fileManager FileManager
var incidentManager IncidentManager
var hookManager HookManager

// HandleCreateIncident handles the create incident web request.
func HandleCreateIncident(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Got Create request")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	incident, pass := convertAdd(r.Body)
	if !pass {
		logManager.LogPrintf("Bad request for create incident %v", incident)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	incident.Type = "Incident"
	incident.State = "open"

	passed := incidentManager.AddIncident(&incident)
	if !passed {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logManager.LogPrintf("Created incident %v\n", incident.Id)
	go hookManager.CallAddedHooks(incident)

	data, err := json.Marshal(incident)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

// HandleIncidentUpdate handles the update incident web request.
func HandleIncidentUpdate(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrint("Got Incident Update")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	incidentId, err := strconv.Atoi(vars["incidentId"])

	if err != nil {
		logManager.LogPrintf("Error converting incidentId %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	update, passed := convertUpdate(r.Body)

	if !passed {
		logManager.LogPrintf("Invalid update for %v\n", incidentId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if incidentManager.UpdateIncident(incidentId, update) {
		go hookManager.CallUpdatedHooks(incidentId, update)
		w.WriteHeader(http.StatusOK)
		return
	}

	logManager.LogPrintf("Incident %v not found\n", incidentId)
	w.WriteHeader(http.StatusNotFound)
}

func convertUpdate(body io.ReadCloser) (IncidentUpdate, bool) {
	decoder := json.NewDecoder(body)

	var update IncidentUpdate
	err := decoder.Decode(&update)

	if err != nil {
		return update, false
	}

	return update, true
}

func convertAdd(body io.ReadCloser) (Incident, bool) {
	decoder := json.NewDecoder(body)

	var inc Incident
	err := decoder.Decode(&inc)

	if err != nil {
		logManager.LogPrintf("Got error when attempting to decode body %v", err)
		return inc, false
	}

	if len(inc.Reporter) == 0 || len(inc.Description) == 0 {
		return inc, false
	}

	return inc, true
}

// HandleGetAttachments handles the get attachment web request.
func HandleGetAttachments(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Getting attachments")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)

	incidentId, err := strconv.Atoi(vars["incidentId"])
	if err != nil {
		logManager.LogPrintf("Error converting incidentId %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := incidentManager.GetIncident(incidentId)
	if !ok {
		logManager.LogPrintf("Got Invalid attachment request for %v.", incidentId)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	attachments, ok := incidentManager.GetAttachments(incidentId)
	if !ok {
		logManager.LogPrintf("Unable to find attachments")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if err := json.NewEncoder(w).Encode(attachments); err != nil {
		panic(err)
	}
}

// HandleUploadAttachment handles the upload attachment web request.
func HandleUploadAttachment(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Got upload attachment request")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	incidentId, err := strconv.Atoi(vars["incidentId"])

	if err != nil {
		logManager.LogPrintf("Error converting incidentId %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := incidentManager.GetIncident(incidentId)
	if !ok {
		logManager.LogPrintf("Got Invalid attachment request for %v.", incidentId)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		logManager.LogPrintln("Unable to get file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	path, ok := fileManager.SaveFile(strconv.Itoa(incidentId), handler.Filename, file)
	if !ok {
		logManager.LogPrintln("Unable to save file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logManager.LogPrintf("Attachment uploaded to %v\n.", path)
	attach := Attachment{handler.Filename, time.Now().Format(time.RFC3339)}
	if !incidentManager.AddAttachment(incidentId, attach) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logManager.LogPrintln("Updated incident with attachment")
	go hookManager.CallAttachedHooks(incidentId, attach)
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	if err := json.NewEncoder(w).Encode(attach); err != nil {
		panic(err)
	}
}

// HandleDownloadAttachment handles the download attachment web request.
func HandleDownloadAttachment(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Got download request.")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)

	incidentId := vars["incidentId"]
	if len(incidentId) <= 0 {
		logManager.LogPrintln("Invalid incident requested.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	attachmentId := vars["attachmentId"]
	if len(attachmentId) <= 0 {
		logManager.LogPrintln("Invalid attachment requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	f, d, passed, callback := fileManager.LoadFile(incidentId, attachmentId)
	if !passed {
		logManager.LogPrintln("File not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer callback()

	http.ServeContent(w, r, d.Name(), d.ModTime(), f)
}

func HandleRemoveAttachment(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Got remove attachment request.")

	vars := mux.Vars(r)

	incident := vars["incidentId"]
	incidentId, err := strconv.Atoi(incident)

	if err != nil {
		logManager.LogPrintf("Error converting incidentId %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := incidentManager.GetIncident(incidentId)

	if !ok {
		logManager.LogPrintf("Got Invalid attachment request for %v.\n", incidentId)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	attachmentId := vars["attachmentId"]
	if len(attachmentId) <= 0 {
		logManager.LogPrintln("Invalid attachment requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	attachments, success := incidentManager.GetAttachments(incidentId)

	if !success {
		logManager.LogFatalf("Unable to get attachments for incident %v.\n", incidentId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	found := false
	for _, v := range attachments {
		if v.FileName == attachmentId {
			found = true
			break
		}
	}

	if !found {
		logManager.LogPrintf("Got invalid attachment request for attachment id %v.\n", attachmentId)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !incidentManager.RemoveAttachment(incidentId, attachmentId) {
		logManager.LogFatalln("Unable to remove attachment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fileManager.DeleteFile(incident, attachmentId)
}

// HandleGetIncident handles the get incident web request.
func HandleGetIncident(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Got incident state request")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)

	incidentId, err := strconv.Atoi(vars["incidentId"])
	if err != nil {
		logManager.LogPrintf("Error converting incidentId %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if val, ok := incidentManager.GetIncident(incidentId); ok {
		logManager.LogPrintf("Got State request for %v.", incidentId)
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(val); err != nil {
			panic(err)
		}

		return
	}

	logManager.LogPrintf("Incident %v not found\n", incidentId)
	w.WriteHeader(http.StatusNotFound)
}

// HandleGetIncidents handles the get incidents web request.
func HandleGetIncidents(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Got incidents request")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")

	filter, passed := convertFilter(r)

	if !passed {
		logManager.LogPrintln("Invalid filter for get request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if filter != nil {
		logManager.LogPrintf("Using filter %+v\n", *filter)
	}

	if val, ok := incidentManager.GetIncidents(filter); ok {
		logManager.LogPrintf("Found %v incidents\n", len(val))

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(val); err != nil {
			logManager.LogPrintln("Unable to encode incidents")
			panic(err)
		}

		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func convertFilter(r *http.Request) (*FilterRequest, bool) {
	param := r.URL.Query()["filter"]
	if param == nil {
		logManager.LogPrintln("Unable to find filter param")
		return nil, true
	}

	/*if err := r.ParseForm(); err != nil {
		logManager.LogPrintf("Unable to parse form: %v\n", err)
		return nil, false
	}*/

	logManager.LogPrintf("got param: %v\n", param[0])

	//data, err := json.Marshal(param[0])

	/*if err != nil {
		logManager.LogPrintf("Unable to marshal param: %v\n", err)
		return nil, false
	}*/

	filter := new(FilterRequest)
	if err := json.Unmarshal([]byte(param[0]), filter); err != nil {
		logManager.LogPrintf("Unable to unmarshal param: %v\n", err)
		return nil, false
	}

	return filter, true
}
