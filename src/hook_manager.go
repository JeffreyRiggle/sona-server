package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// HookManager provides the ability to make web endpoint calls.
// The AddedWebHooks are the endpoints to call in CallAddedHooks.
// The UpdatedWebHooks are the endpoints to call in CallUpdatedHooks.
// The AttachedWebHooks are the endpoints to call in CallAttachedWebHooks.
type HookManager struct {
	AddedWebHooks    []WebHook
	UpdatedWebHooks  []WebHook
	AttachedWebHooks []WebHook
}

// CallAddedHooks will call all defined added endpoints.
// During this process it will subsitute any nessicary data.
func (manager HookManager) CallAddedHooks(incident Incident) {
	logManager.LogPrintln("Calling added hooks")
	for _, hook := range manager.AddedWebHooks {
		go fireHook(hook, preformAddedSubsitutions(hook, incident))
	}
}

// CallUpdatedHooks will call all defined updated endpoints.
// During this process it will subsitute any nessicary data.
func (manager HookManager) CallUpdatedHooks(incidentID int, incident IncidentUpdate) {
	logManager.LogPrintln("Calling updated hooks")
	for _, hook := range manager.UpdatedWebHooks {
		go fireHook(hook, preformUpdateSubsitutions(hook, incidentID, incident))
	}
}

// CallAttachedHooks will call all defined attached endpoints.
// During this process it will subsitute any nessicary data.
func (manager HookManager) CallAttachedHooks(incidentID int, attachment Attachment) {
	logManager.LogPrintln("Calling attached hooks")
	for _, hook := range manager.AttachedWebHooks {
		go fireHook(hook, preformAttachSubsitutions(hook, incidentID, attachment))
	}
}

func preformAddedSubsitutions(hook WebHook, incident Incident) *bytes.Buffer {
	var bod = make(map[string]string, 0)

	for _, item := range hook.Body.Items {
		if item.Substitute {
			bod[item.Key] = preformAddSubstitutionImpl(item.Value, incident)
		} else {
			bod[item.Key] = item.Value
		}
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(bod)
	return b
}

func preformAddSubstitutionImpl(key string, incident Incident) string {
	var cRegEx = regexp.MustCompile("\\{\\{([^\\}\\}]*)\\}\\}")
	match := cRegEx.FindAllStringSubmatch(key, -1)

	if len(match) <= 0 {
		return getIncidentPropertyValue(key, incident)
	}

	var retVal = key
	for i := 0; i < len(match); i++ {
		var replaceRegEx = regexp.MustCompile(match[i][0])
		retVal = replaceRegEx.ReplaceAllString(retVal, getIncidentPropertyValue(match[i][1], incident))
	}

	return retVal
}

func preformUpdateSubsitutions(hook WebHook, incidentID int, incident IncidentUpdate) *bytes.Buffer {
	var bod = make(map[string]string, 0)

	for _, item := range hook.Body.Items {
		if item.Substitute {
			bod[item.Key] = preformUpdateSubstitutionImpl(item.Value, incidentID, incident)
		} else {
			bod[item.Key] = item.Value
		}
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(bod)
	return b
}

func preformUpdateSubstitutionImpl(key string, incidentID int, incident IncidentUpdate) string {
	var cRegEx = regexp.MustCompile("\\{\\{([^\\}\\}]*)\\}\\}")
	match := cRegEx.FindAllStringSubmatch(key, -1)

	if len(match) <= 0 {
		return getUpdateSubstitutionValue(key, incidentID, incident)
	}

	var retVal = key
	for i := 0; i < len(match); i++ {
		var replaceRegEx = regexp.MustCompile(match[i][0])
		retVal = replaceRegEx.ReplaceAllString(retVal, getUpdateSubstitutionValue(match[i][1], incidentID, incident))
	}

	return retVal
}

func getUpdateSubstitutionValue(key string, incidentID int, incident IncidentUpdate) string {
	if key == "id" {
		return strconv.Itoa(incidentID)
	}
	if key == "reporter" {
		return incident.Reporter
	}
	if key == "description" {
		return incident.Description
	}
	if key == "state" {
		return incident.State
	}

	if val, ok := incident.Attributes[key]; ok {
		return val
	}

	return ""
}

func preformAttachSubsitutions(hook WebHook, incidentID int, attachment Attachment) *bytes.Buffer {
	var bod = make(map[string]string, 0)

	for _, item := range hook.Body.Items {
		if item.Substitute {
			bod[item.Key] = preformAttachSubstitutionImpl(item.Value, incidentID, attachment)
		} else {
			bod[item.Key] = item.Value
		}
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(bod)
	return b
}

func preformAttachSubstitutionImpl(key string, incidentID int, attachment Attachment) string {
	var cRegEx = regexp.MustCompile("\\{\\{([^\\}\\}]*)\\}\\}")
	match := cRegEx.FindAllStringSubmatch(key, -1)

	if len(match) <= 0 {
		return getAttachSubstitutionValue(key, incidentID, attachment)
	}

	var retVal = key
	for i := 0; i < len(match); i++ {
		var replaceRegEx = regexp.MustCompile(match[i][0])
		retVal = replaceRegEx.ReplaceAllString(retVal, getAttachSubstitutionValue(match[i][1], incidentID, attachment))
	}

	return retVal
}

func getAttachSubstitutionValue(key string, incidentID int, attachment Attachment) string {
	if key == "id" {
		return strconv.Itoa(incidentID)
	}
	if key == "attachment" {
		b, err := json.Marshal(attachment)
		if err != nil {
			logManager.LogFatal("Unable to create attachment json")
			return ""
		}
		return string(b)
	}
	if key == "filename" {
		return attachment.FileName
	}

	return ""
}

func fireHook(hook WebHook, body *bytes.Buffer) {
	client := http.Client{
		Timeout: time.Second * 5,
	}

	logManager.LogPrintf("Calling %v %v with body %v\n", hook.Method, hook.URL, body)
	req, err := http.NewRequest(hook.Method, hook.URL, body)

	if err != nil {
		logManager.LogPrintf("Failure creating webhook %v: %v\n", hook, err)
	}

	_, err2 := client.Do(req)

	if err2 != nil {
		logManager.LogPrintf("Unable failure using webhook %v: %v\n", hook, err)
	}
}
