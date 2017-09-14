package main

// IncidentManager defines a minimal implementation required for managing incidents.
// AddIncident should add an incident to the manager.
// GetIncident should return the requested incident and return a false if the incident does not exist
// GetIncidents should return all managed incidents.
// Update incident should update the underlying incident with new data.
// AddAttachments should update the association between an incident and an attachment.
// GetAttachments should get all attachments associated with an incident.
type IncidentManager interface {
	AddIncident(incident *Incident) bool
	GetIncident(incidentId int) (Incident, bool)
	GetIncidents() ([]Incident, bool)
	UpdateIncident(id int, incident IncidentUpdate) bool
	AddAttachment(incidentId int, attachment Attachment) bool
	GetAttachments(incidentId int) ([]Attachment, bool)
}
