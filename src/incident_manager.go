package main

// IncidentManager defines a minimal implementation required for managing incidents.
// AddIncident should add an incident to the manager.
// GetIncident should return the requested incident and return a false if the incident does not exist
// GetIncidents should return all managed incidents.
// Update incident should update the underlying incident with new data.
// AddAttachments should update the association between an incident and an attachment.
// GetAttachments should get all attachments associated with an incident.
// RemoveAttachment will find and remove an attachment associated with an incident.
// CleanUp will do any required cleanup actions on the incident manager.
type IncidentManager interface {
	AddIncident(incident *Incident) bool
	GetIncident(incidentId int) (Incident, bool)
	GetIncidents(filter *FilterRequest) ([]Incident, bool)
	UpdateIncident(id int, incident IncidentUpdate) bool
	AddAttachment(incidentId int, attachment Attachment) bool
	GetAttachments(incidentId int) ([]Attachment, bool)
	RemoveAttachment(incidentId int, fileName string) bool
	CleanUp()
}
