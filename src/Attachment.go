package main

// Attachment defines a file attached to an incident.
type Attachment struct {
	FileName string `json:"filename"` // The file name.
	Time     string `json:"time"`     // The time the file was attached.
}
