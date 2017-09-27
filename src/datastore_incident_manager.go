package main

import (
	"strconv"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type DataStoreIncidentManager struct {
	Context     *context.Context
	Connection  *datastore.Client
	LastKnownId int64
}

type DataStoreIncident struct {
	Type        string
	Id          int64
	Description string
	Reporter    string
	State       string
}

func CreateDataStoreClient(projectName string, authFile string) (*context.Context, *datastore.Client) {
	context := context.Background()
	client, err := datastore.NewClient(context, projectName, option.WithServiceAccountFile(authFile))

	if err != nil {
		panic(err)
	}

	return &context, client
}

func (manager DataStoreIncidentManager) AddIncident(incident *Incident) bool {
	id := manager.getNextId()
	incident.Id = id

	dsInc := DataStoreIncident{incident.Type, incident.Id, incident.Description, incident.Reporter, incident.State}
	taskKey := datastore.NameKey("incidents", strconv.FormatInt(id, 10), nil)

	logManager.LogPrintln("Attempting to put incident into database")

	if _, err := manager.Connection.Put(*manager.Context, taskKey, &dsInc); err != nil {
		logManager.LogPrintf("Unable to put incident %v\n", err)
		return false
	}

	manager.LastKnownId = id
	logManager.LogPrintln("Incident created in database")
	return true
}

func (manager DataStoreIncidentManager) getNextId() int64 {
	logManager.LogPrintf("Getting next id last known id %v\n", manager.LastKnownId)

	q := datastore.NewQuery("incidents").Filter("Id >", manager.LastKnownId)
	iter := manager.Connection.Run(*manager.Context, q)

	max := manager.LastKnownId
	for {
		var incident Incident
		_, err := iter.Next(&incident)

		if err == iterator.Done {
			break
		}

		if err != nil {
			logManager.LogPrintf("Got error when attempting to get ids %v\n", err)
			break
		}

		if incident.Id > max {
			max = incident.Id
		}
	}

	logManager.LogPrintf("Got new max id %v\n", max+1)
	return max + 1
}

func (manager DataStoreIncidentManager) GetIncident(incidentId int) (Incident, bool) {
	q := datastore.NewQuery("incidents").Filter("Id ==", incidentId)
	iter := manager.Connection.Run(*manager.Context, q)

	var incident Incident
	for {
		_, err := iter.Next(&incident)

		if err == iterator.Done {
			break
		}
	}

	return incident, true
}

func (manager DataStoreIncidentManager) GetIncidents() ([]Incident, bool) {
	retVal := make([]Incident, 0)

	q := datastore.NewQuery("incidents")
	iter := manager.Connection.Run(*manager.Context, q)

	for {
		var incident Incident
		_, err := iter.Next(&incident)

		if err == iterator.Done {
			break
		}

		retVal = append(retVal, incident)
	}

	return retVal, true
}

func (manager DataStoreIncidentManager) UpdateIncident(id int, incident IncidentUpdate) bool {
	return true
}

func (manager DataStoreIncidentManager) AddAttachment(incidentId int, attachment Attachment) bool {
	return true
}

func (manager DataStoreIncidentManager) GetAttachments(incidentId int) ([]Attachment, bool) {
	return make([]Attachment, 0), true
}

func (manager DataStoreIncidentManager) RemoveAttachment(incidentId int, fileName string) bool {
	return true
}

func (manager DataStoreIncidentManager) CleanUp() {

}
