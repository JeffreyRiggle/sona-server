package main

// RuntimeIncidentManager manages incidents in the applications runtime.
// These incidents will no longer be available after the application shuts down.
type RuntimeIncidentManager struct {
	Incidents   map[int]*Incident    // The incidents created.
	Attachments map[int][]Attachment // The attachments and the incident association.
}

// AddIncident adds an incident to the runtimes incident collection.
func (manager RuntimeIncidentManager) AddIncident(incident *Incident) bool {
	var id = len(manager.Incidents)
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
	if val, ok := manager.Incidents[incidentId]; ok {
		return *val, true
	}

	return Incident{}, false
}

// GetIncidents will get all incidents out of the runtimes incident collection.
func (manager RuntimeIncidentManager) GetIncidents() ([]Incident, bool) {
	retVal := make([]Incident, 0)

	for _, v := range manager.Incidents {
		retVal = append(retVal, *v)
	}

	return retVal, true
}

// UpdateIncident will update a given incident in the runtime.
func (manager RuntimeIncidentManager) UpdateIncident(id int, update IncidentUpdate) bool {
	if val, ok := manager.Incidents[id]; ok {
		updateIncident(val, update)
		return true
	}

	return false
}

// AddAttachment will create an association between an attachment and an incident in the runtime.
func (manager RuntimeIncidentManager) AddAttachment(incidentId int, attachment Attachment) bool {
	if _, ok := manager.Incidents[incidentId]; ok {
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
