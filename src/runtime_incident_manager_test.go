package main

import "testing"

func TestAddIncident(t *testing.T) {
	var manager = RuntimeIncidentManager{make(map[int64]*Incident, 0), make(map[int][]Attachment, 0)}
	manager.AddIncident(new(Incident))

	if len(manager.Incidents) != 1 {
		t.Error(
			"For", manager,
			"expected", 1,
			"got", len(manager.Incidents))
	}

	if manager.Incidents[0].Id != 0 {
		t.Error(
			"For", manager.Incidents[0],
			"expected", 0,
			"got", manager.Incidents[0].Id)
	}

	manager.AddIncident(new(Incident))

	if len(manager.Incidents) != 2 {
		t.Error(
			"For", manager,
			"expected", 1,
			"got", len(manager.Incidents))
	}

	if manager.Incidents[1].Id != 1 {
		t.Error(
			"For", manager.Incidents[1],
			"expected", 1,
			"got", manager.Incidents[1].Id)
	}
}

func TestGetIncident(t *testing.T) {
	var manager = RuntimeIncidentManager{make(map[int64]*Incident, 0), make(map[int][]Attachment, 0)}
	var incident = Incident{"Incident", 0, "Some Description", "Someone", "Open", make(map[string]string, 0)}
	manager.AddIncident(&incident)

	retVal, pass := manager.GetIncident(0)

	if !pass {
		t.Error(
			"For", pass,
			"expected", true,
			"got", pass)
	}

	if retVal.Id != 0 {
		t.Error(
			"For", retVal,
			"expected", 0,
			"got", retVal.Id)
	}

	if retVal.Description != "Some Description" {
		t.Error(
			"For", retVal,
			"expected", "Some Description",
			"got", retVal.Description)
	}

	if retVal.Reporter != "Someone" {
		t.Error(
			"For", retVal,
			"expected", "Someone",
			"got", retVal.Reporter)
	}

	if retVal.State != "Open" {
		t.Error(
			"For", retVal,
			"expected", "Open",
			"got", retVal.State)
	}

	if len(retVal.Attributes) != 0 {
		t.Error(
			"For", retVal,
			"expected", 0,
			"got", len(retVal.Attributes))
	}
}

func TestGetInvalidIncident(t *testing.T) {
	var manager = RuntimeIncidentManager{make(map[int64]*Incident, 0), make(map[int][]Attachment, 0)}
	var incident = Incident{"Incident", 0, "Some Description", "Someone", "Open", make(map[string]string, 0)}
	manager.AddIncident(&incident)

	_, pass := manager.GetIncident(1)

	if pass {
		t.Error(
			"For", pass,
			"expected", false,
			"got", pass)
	}
}

func TestGetIncidents(t *testing.T) {
	var manager = RuntimeIncidentManager{make(map[int64]*Incident, 0), make(map[int][]Attachment, 0)}
	var incident1 = Incident{"Incident", 0, "Some Description", "Someone", "Open", make(map[string]string, 0)}
	var incident2 = Incident{"Incident", 0, "Some Description", "Someone", "Open", make(map[string]string, 0)}
	manager.AddIncident(&incident1)
	manager.AddIncident(&incident2)

	retVal, pass := manager.GetIncidents(nil)

	if !pass {
		t.Error(
			"For", pass,
			"expected", true,
			"got", pass)
	}

	if len(retVal) != 2 {
		t.Error(
			"For", retVal,
			"expected", 2,
			"got", len(retVal))
	}
}

func TestUpdateIncident(t *testing.T) {
	var manager = RuntimeIncidentManager{make(map[int64]*Incident, 0), make(map[int][]Attachment, 0)}
	var incident = Incident{"Incident", 0, "Some Description", "Someone", "Open", make(map[string]string, 0)}
	manager.AddIncident(&incident)
	manager.UpdateIncident(0, IncidentUpdate{"New State", "New Description", "", nil})

	retVal, pass := manager.GetIncident(0)

	if !pass {
		t.Error(
			"For", pass,
			"expected", true,
			"got", pass)
	}

	if retVal.Id != 0 {
		t.Error(
			"For", retVal,
			"expected", 0,
			"got", retVal.Id)
	}

	if retVal.Description != "New Description" {
		t.Error(
			"For", retVal,
			"expected", "New Description",
			"got", retVal.Description)
	}

	if retVal.Reporter != "Someone" {
		t.Error(
			"For", retVal,
			"expected", "Someone",
			"got", retVal.Reporter)
	}

	if retVal.State != "New State" {
		t.Error(
			"For", retVal,
			"expected", "New State",
			"got", retVal.State)
	}

	if len(retVal.Attributes) != 0 {
		t.Error(
			"For", retVal,
			"expected", 0,
			"got", len(retVal.Attributes))
	}
}

func TestAddAttachment(t *testing.T) {
	var manager = RuntimeIncidentManager{make(map[int64]*Incident, 0), make(map[int][]Attachment, 0)}
	var incident = Incident{"Incident", 0, "Some Description", "Someone", "Open", make(map[string]string, 0)}
	manager.AddIncident(&incident)

	var attach = Attachment{"testfile.jpg", "2009-11-10T23:00:00Z"}
	pass := manager.AddAttachment(0, attach)

	if !pass {
		t.Error(
			"For", pass,
			"expected", true,
			"got", pass)
	}

	if len(manager.Attachments) != 1 {
		t.Error(
			"For", manager.Attachments,
			"expected", 1,
			"got", len(manager.Attachments))
	}
}

func TestAddAttachmentToInvalidIncident(t *testing.T) {
	var manager = RuntimeIncidentManager{make(map[int64]*Incident, 0), make(map[int][]Attachment, 0)}
	var incident = Incident{"Incident", 0, "Some Description", "Someone", "Open", make(map[string]string, 0)}
	manager.AddIncident(&incident)

	var attach = Attachment{"testfile.jpg", "2009-11-10T23:00:00Z"}
	pass := manager.AddAttachment(1, attach)

	if pass {
		t.Error(
			"For", pass,
			"expected", false,
			"got", pass)
	}

	if len(manager.Attachments) != 0 {
		t.Error(
			"For", manager.Attachments,
			"expected", 0,
			"got", len(manager.Attachments))
	}
}

func TestGetAttachments(t *testing.T) {
	var manager = RuntimeIncidentManager{make(map[int64]*Incident, 0), make(map[int][]Attachment, 0)}
	var incident = Incident{"Incident", 0, "Some Description", "Someone", "Open", make(map[string]string, 0)}
	manager.AddIncident(&incident)

	var attach1 = Attachment{"testfile.jpg", "2009-11-10T23:00:00Z"}
	var attach2 = Attachment{"testfile2.jpg", "2009-10-10T23:00:00Z"}
	pass1 := manager.AddAttachment(0, attach1)
	pass2 := manager.AddAttachment(0, attach2)

	if !pass1 {
		t.Error(
			"For", pass1,
			"expected", true,
			"got", pass1)
	}

	if !pass2 {
		t.Error(
			"For", pass2,
			"expected", true,
			"got", pass2)
	}

	retVal, pass3 := manager.GetAttachments(0)

	if !pass3 {
		t.Error(
			"For", pass3,
			"expected", true,
			"got", pass3)
	}

	if len(retVal) != 2 {
		t.Error(
			"For", retVal,
			"expected", 2,
			"got", len(retVal))
	}

	if retVal[0].FileName != "testfile.jpg" {
		t.Error(
			"For", retVal[0],
			"expected", "testfile.jpg",
			"got", retVal[0].FileName)
	}

	if retVal[1].FileName != "testfile2.jpg" {
		t.Error(
			"For", retVal[1],
			"expected", "testfile2.jpg",
			"got", retVal[1].FileName)
	}
}

func TestRemoveAttribute(t *testing.T) {
	var manager = RuntimeIncidentManager{make(map[int64]*Incident, 0), make(map[int][]Attachment, 0)}
	var incident = Incident{"Incident", 0, "Some Description", "Someone", "Open", make(map[string]string, 0)}
	manager.AddIncident(&incident)
	var attributes = make(map[string]string, 0)
	attributes["Test"] = "val"
	attributes["Test2"] = "var"
	manager.UpdateIncident(0, IncidentUpdate{"New State", "New Description", "", attributes})

	retVal, pass := manager.GetIncident(0)

	if !pass {
		t.Error(
			"For", pass,
			"expected", true,
			"got", pass)
	}

	if len(retVal.Attributes) != 2 {
		t.Error(
			"For", retVal,
			"expected", 2,
			"got", len(retVal.Attributes))
	}

	var attributes2 = make(map[string]string, 0)
	attributes2["Test"] = "val"
	manager.UpdateIncident(0, IncidentUpdate{"New State", "New Description", "", attributes2})

	retVal2, pass2 := manager.GetIncident(0)

	if !pass2 {
		t.Error(
			"For", pass2,
			"Expected", true,
			"got", pass2)
	}

	if len(retVal2.Attributes) != 1 {
		t.Error(
			"For", retVal2,
			"expected", 1,
			"got", len(retVal2.Attributes))
	}
}

func TestRemoveAttachment(t *testing.T) {
	var manager = RuntimeIncidentManager{make(map[int64]*Incident, 0), make(map[int][]Attachment, 0)}
	var incident = Incident{"Incident", 0, "Some Description", "Someone", "Open", make(map[string]string, 0)}
	manager.AddIncident(&incident)

	var attach1 = Attachment{"testfile.jpg", "2009-11-10T23:00:00Z"}
	var attach2 = Attachment{"testfile2.jpg", "2009-10-10T23:00:00Z"}
	var attach3 = Attachment{"testfile3.jpg", "2009-12-10T23:00:00Z"}
	pass1 := manager.AddAttachment(0, attach1)
	pass2 := manager.AddAttachment(0, attach2)
	pass3 := manager.AddAttachment(0, attach3)

	if !pass1 {
		t.Error(
			"For", pass1,
			"expected", true,
			"got", pass1)
	}

	if !pass2 {
		t.Error(
			"For", pass2,
			"expected", true,
			"got", pass2)
	}

	if !pass3 {
		t.Error(
			"For", pass3,
			"expected", true,
			"got", pass3)
	}

	pass4 := manager.RemoveAttachment(0, "testfile2.jpg")

	if !pass4 {
		t.Error(
			"For", pass4,
			"expected", true,
			"got", pass4)
	}

	retVal, pass5 := manager.GetAttachments(0)

	if !pass5 {
		t.Error(
			"For", pass3,
			"expected", true,
			"got", pass3)
	}

	if len(retVal) != 2 {
		t.Error(
			"For", retVal,
			"expected", 2,
			"got", len(retVal))
	}

	if retVal[0].FileName != "testfile.jpg" {
		t.Error(
			"For", retVal[0],
			"expected", "testfile.jpg",
			"got", retVal[0].FileName)
	}

	if retVal[1].FileName != "testfile3.jpg" {
		t.Error(
			"For", retVal[1],
			"expected", "testfile3.jpg",
			"got", retVal[1].FileName)
	}
}
