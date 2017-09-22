package main

// Incident defines the basic item for managing and tracking issues.
type Incident struct {
	Type        string            `json:"type"`        // The type of incident.
	Id          int64             `json:"id"`          // The unique identifier for the incident.
	Description string            `json:"description"` // The description of the incident.
	Reporter    string            `json:"reporter"`    // The reporter of the incident.
	State       string            `json:"state"`       // The current state of the incident.
	Attributes  map[string]string `json:"attributes"`  // The attributes associated with the incident.
}

// IncidentUpdate defines a set of updates to apply to an underyling incident.
type IncidentUpdate struct {
	State       string            `json:"state"`       // The new state for the incident.
	Description string            `json:"description"` // The new description of the incident.
	Reporter    string            `json:"reporter"`    // The new reporter of the incident.
	Attributes  map[string]string `json:"attributes"`  // The new attributes to associate with the incident.
}

func updateIncident(original *Incident, updated IncidentUpdate) bool {
	changed := false
	if len(updated.State) > 0 {
		original.State = updated.State
		changed = true
	}

	if len(updated.Reporter) > 0 {
		original.Reporter = updated.Reporter
		changed = true
	}

	if len(updated.Description) > 0 {
		original.Description = updated.Description
		changed = true
	}

	if updated.Attributes != nil {
		original.Attributes = updated.Attributes
		changed = true
	}

	return changed
}
