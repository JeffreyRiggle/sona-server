package main

import (
	"sort"
	"strings"
)

// RuntimeIncidentManager manages incidents in the applications runtime.
// These incidents will no longer be available after the application shuts down.
type RuntimeIncidentManager struct {
	Incidents   map[int64]*Incident  // The incidents created.
	Attachments map[int][]Attachment // The attachments and the incident association.
}

// AddIncident adds an incident to the runtimes incident collection.
func (manager RuntimeIncidentManager) AddIncident(incident *Incident) bool {
	var id = int64(len(manager.Incidents))
	incident.Id = id
	if incident.Attributes == nil {
		incident.Attributes = make(map[string]string, 0)
	}

	manager.Incidents[id] = incident
	return true
}

// GetIncident attempts to get an incident out of the runtimes incident collection.
// If an incident is not found a false will be returned.
func (manager RuntimeIncidentManager) GetIncident(incidentId int) (Incident, bool) {
	if val, ok := manager.Incidents[int64(incidentId)]; ok {
		return *val, true
	}

	return Incident{}, false
}

// GetIncidents will get all incidents out of the runtimes incident collection.
func (manager RuntimeIncidentManager) GetIncidents(filter *FilterRequest) ([]Incident, bool) {
	retVal := make([]Incident, 0)

	for _, v := range manager.Incidents {
		if incidentInFilterRequest(*v, filter) {
			retVal = append(retVal, *v)
		}
	}

	sort.Slice(retVal, func(i, j int) bool {
		return retVal[i].Id < retVal[j].Id
	})

	return retVal, true
}

func incidentInFilterRequest(incident Incident, filter *FilterRequest) bool {
	if filter == nil {
		return true
	}

	if isOrRequest(filter) {
		for _, k := range filter.Filters {
			if incidentInComplexFilter(incident, k) {
				return true
			}
		}

		return false
	}

	for _, k := range filter.Filters {
		if !incidentInComplexFilter(incident, k) {
			return false
		}
	}

	return true
}

func incidentInComplexFilter(incident Incident, filter ComplexFilter) bool {
	logManager.LogPrintf("Processing complex filter %v\n", filter)

	if filter.Children != nil {
		logManager.LogPrintf("Processing complex filter with children %v\n", filter.Children)

		if isOrFilter(filter) {
			for _, v := range filter.Children {
				if incidentInComplexFilter(incident, *v) {
					return true
				}
			}

			return false
		}

		for _, v := range filter.Children {
			if !incidentInComplexFilter(incident, *v) {
				return false
			}
		}

		return true
	}

	if isOrFilter(filter) {
		for _, v := range filter.Filter {
			if incidentInFilter(incident, v) {
				return true
			}
		}

		return false
	}

	for _, v := range filter.Filter {
		if !incidentInFilter(incident, v) {
			return false
		}
	}

	return true
}

func incidentInFilter(incident Incident, filter Filter) bool {
	val := getIncidentPropertyValue(filter.Property, incident)
	if isEqualsComparision(filter) {
		return strings.EqualFold(filter.Value, val)
	}

	if isContainsComparision(filter) {
		return strings.Contains(strings.ToLower(val), strings.ToLower(filter.Value))
	}

	if isNotEqualsComparision(filter) {
		return !strings.EqualFold(filter.Value, val)
	}

	return false
}

// UpdateIncident will update a given incident in the runtime.
func (manager RuntimeIncidentManager) UpdateIncident(id int, update IncidentUpdate) bool {
	if val, ok := manager.Incidents[int64(id)]; ok {
		updateIncident(val, update)
		return true
	}

	return false
}

// AddAttachment will create an association between an attachment and an incident in the runtime.
func (manager RuntimeIncidentManager) AddAttachment(incidentId int, attachment Attachment) bool {
	if _, ok := manager.Incidents[int64(incidentId)]; ok {
		manager.Attachments[incidentId] = append(manager.Attachments[incidentId], attachment)
		return true
	}

	return false
}

// GetAttachments will find all attachments associated with an incident in the runtime.
func (manager RuntimeIncidentManager) GetAttachments(incidentId int) ([]Attachment, bool) {
	if val, ok := manager.Attachments[incidentId]; ok {
		return val, true
	}

	attachments := make([]Attachment, 0)
	manager.Attachments[incidentId] = attachments
	return attachments, true
}

// RemoveAttachment will find and remove an attachment associated with an incident.
func (manager RuntimeIncidentManager) RemoveAttachment(incidentId int, fileName string) bool {
	val, ok := manager.Attachments[incidentId]

	if !ok {
		return false
	}

	for i, v := range val {
		if v.FileName == fileName {
			manager.Attachments[incidentId] = append(val[:i], val[i+1:]...)
			return true
		}
	}

	return false
}

// CleanUp will do any required cleanup actions on the incident manager.
func (manager RuntimeIncidentManager) CleanUp() {
	// No op
}
