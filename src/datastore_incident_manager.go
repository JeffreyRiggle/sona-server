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
	Attributes  []DataStoreIncidentAttribute
}

type DataStoreIncidentAttribute struct {
	Name  string
	Value string
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

	dsInc := convertFromIncident(incident)
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

func convertFromIncident(incident *Incident) DataStoreIncident {
	attributes := make([]DataStoreIncidentAttribute, 0)
	for k, v := range incident.Attributes {
		attributes = append(attributes, DataStoreIncidentAttribute{Name: k, Value: v})
	}

	return DataStoreIncident{incident.Type, incident.Id, incident.Description, incident.Reporter, incident.State, attributes}
}

func convertToIncident(incident *DataStoreIncident) Incident {
	retVal := Incident{incident.Type, incident.Id, incident.Description, incident.Reporter, incident.State, make(map[string]string, 0)}
	for _, v := range incident.Attributes {
		retVal.Attributes[v.Name] = v.Value
	}

	return retVal
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
	q := datastore.NewQuery("incidents").Filter("Id =", incidentId)
	iter := manager.Connection.Run(*manager.Context, q)

	var incident DataStoreIncident
	for {
		_, err := iter.Next(&incident)

		if err == iterator.Done {
			break
		}

		if err != nil {
			logManager.LogPrintf("Got error when attempting to get incident %v\n", err)
			break
		}
	}

	return convertToIncident(&incident), true
}

func (manager DataStoreIncidentManager) GetIncidents() ([]Incident, bool) {
	retVal := make([]Incident, 0)

	q := datastore.NewQuery("incidents")
	iter := manager.Connection.Run(*manager.Context, q)

	for {
		var incident DataStoreIncident
		_, err := iter.Next(&incident)

		if err == iterator.Done {
			break
		}

		if err != nil {
			logManager.LogPrintf("Got error when attempting to get incident %v\n", err)
			break
		}

		retVal = append(retVal, convertToIncident(&incident))
	}

	return retVal, true
}

func (manager DataStoreIncidentManager) UpdateIncident(id int, incident IncidentUpdate) bool {
	logManager.LogPrintf("Got incident update request for %v\n", id)
	inc, found := manager.GetIncident(id)

	if !found {
		return false
	}

	if !updateIncident(&inc, incident) {
		return true
	}

	dsInc := convertFromIncident(&inc)
	taskKey := datastore.NameKey("incidents", strconv.FormatInt(inc.Id, 10), nil)

	logManager.LogPrintf("Attempting to update incident %v in database. %v\n", inc.Id, dsInc)

	if _, err := manager.Connection.Put(*manager.Context, taskKey, &dsInc); err != nil {
		logManager.LogPrintf("Unable to update incident %v\n", err)
		return false
	}

	logManager.LogPrintf("Updated incident %v in database\n", inc.Id)

	return true
}

func (manager DataStoreIncidentManager) AddAttachment(incidentId int, attachment Attachment) bool {
	parentKey := datastore.NameKey("incidents", strconv.Itoa(incidentId), nil)
	taskKey := datastore.NameKey("incidentattachments", attachment.FileName, parentKey)

	logManager.LogPrintln("Attempting to put incident attachment into database")

	if _, err := manager.Connection.Put(*manager.Context, taskKey, &attachment); err != nil {
		logManager.LogPrintf("Unable to put incident attachment %v\n", err)
		return false
	}

	logManager.LogPrintln("Incident attachment created in database")
	return true
}

func (manager DataStoreIncidentManager) GetAttachments(incidentId int) ([]Attachment, bool) {
	retVal := make([]Attachment, 0)
	q := datastore.NewQuery("incidentattachments").Ancestor(datastore.NameKey("incidents", strconv.Itoa(incidentId), nil))
	iter := manager.Connection.Run(*manager.Context, q)

	var att Attachment
	for {
		_, err := iter.Next(&att)

		if err == iterator.Done {
			break
		}

		if err != nil {
			logManager.LogPrintf("Got error when attempting to get incident %v\n", err)
			break
		}

		retVal = append(retVal, att)
	}

	return retVal, true
}

func (manager DataStoreIncidentManager) RemoveAttachment(incidentId int, fileName string) bool {
	parentKey := datastore.NameKey("incidents", strconv.Itoa(incidentId), nil)
	taskKey := datastore.NameKey("incidentattachments", fileName, parentKey)

	logManager.LogPrintln("Attempting to remove incident attachment from database")

	err := manager.Connection.Delete(*manager.Context, taskKey)

	if err != nil {
		logManager.LogPrintf("Unable to delete attachment %v", err)
		return false
	}

	return true
}

func (manager DataStoreIncidentManager) CleanUp() {
	if manager.Connection != nil {
		manager.Connection.Close()
	}
}
