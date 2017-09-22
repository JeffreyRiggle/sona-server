package main

import (
	"database/sql"
	"fmt"
)

type MySQLManager struct {
	Connection *sql.DB
}

func (manager MySQLManager) Initialize() {
	if !manager.hasTable("Incidents") {
		logManager.LogPrintln("Unable to find incident table creating now")
		manager.createIncidentTable()
	}

	if !manager.hasTable("IncidentAttributes") {
		logManager.LogPrintln("Unable to find incident attribute table creating now")
		manager.createAttributeTable()
	}

	if !manager.hasTable("IncidentAttachments") {
		logManager.LogPrintln("Unable to find attachment table creating now")
		manager.createAttachmentTable()
	}
}

func (manager MySQLManager) hasTable(tableName string) bool {
	rows, err := manager.Connection.Query(fmt.Sprintf("SHOW TABLES LIKE '%v'", tableName))

	if err != nil {
		logManager.LogPrintf("Got error %v\n", err)
		return false
	}
	defer rows.Close()

	logManager.LogPrintf("Got rows %v\n", rows)
	for rows.Next() {
		return true
	}

	return false
}

func (manager MySQLManager) createIncidentTable() {
	stmt, err := manager.Connection.Prepare("CREATE TABLE Incidents (" +
		"Id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY, " +
		"Type VARCHAR(255), " +
		"Description VARCHAR(1048), " +
		"Reporter VARCHAR(255), " +
		"State VARCHAR(255))")

	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec()
	if err != nil {
		panic(err)
	}

	logManager.LogPrintln("Created Incident Table: %v\n", res)
}

func (manager MySQLManager) createAttributeTable() {
	stmt, err := manager.Connection.Prepare("CREATE TABLE IncidentAttributes (" +
		"IncidentId INT UNSIGNED NOT NULL, " +
		"AttributeName VARCHAR(255), " +
		"AttributeValue VARCHAR(255)," +
		"PRIMARY KEY(IncidentId, AttributeName), " +
		"FOREIGN KEY (IncidentId) " +
		"	REFERENCES Incidents(Id))")

	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec()
	if err != nil {
		panic(err)
	}

	logManager.LogPrintf("Created Attribute Table: %v\n", res)
}

func (manager MySQLManager) createAttachmentTable() {
	stmt, err := manager.Connection.Prepare("CREATE TABLE IncidentAttachments (" +
		"IncidentId INT UNSIGNED NOT NULL, " +
		"FileName VARCHAR(255), " +
		"TimeStampString VARCHAR(255), " +
		"PRIMARY KEY(IncidentId, FileName), " +
		"FOREIGN KEY (IncidentId) " +
		"	REFERENCES Incidents(Id))")

	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec()
	if err != nil {
		panic(err)
	}

	logManager.LogPrintf("Created Attribute Table: %v\n", res)
}

func (manager MySQLManager) AddIncident(incident *Incident) bool {
	stmt, err := manager.Connection.Prepare("INSERT INTO incidents (Type, Description, Reporter, State) " +
		"VALUES (?, ?, ?, ?);")
	if err != nil {
		logManager.LogPrintf("Error occurred when preparing add %v", err)
		return false
	}

	res, err := stmt.Exec(incident.Type, incident.Description, incident.Reporter, incident.State)

	if err != nil {
		logManager.LogPrintf("Error occurred when executing add %v", err)
		return false
	}

	id, err := res.LastInsertId()

	if err != nil {
		logManager.LogPrintf("Unable to get last inserted id %v\n", err)
		return false
	}

	incident.Id = id
	logManager.LogPrintf("Created new incident: %v\n", incident)
	return true
}

func (manager MySQLManager) GetIncident(incidentId int) (Incident, bool) {
	retVal := Incident{"", 0, "", "", "", nil}
	var (
		id           int64
		incidenttype string
		description  string
		reporter     string
		state        string
		attname      sql.NullString
		attvalue     sql.NullString
	)

	rows, err := manager.Connection.Query("SELECT Id, Type, Description, Reporter, State, AttributeName, AttributeValue "+
		"FROM incidents LEFT JOIN incidentattributes "+
		"ON IncidentId = Id "+
		"WHERE Id = ?", incidentId)

	if err != nil {
		logManager.LogPrintf("Error occurred when preparing get %v\n", err)
		return retVal, false
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &incidenttype, &description, &reporter, &state, &attname, &attvalue)
		if err != nil {
			logManager.LogPrintln(err)
		}

		if retVal.Id == 0 {
			retVal = Incident{incidenttype, id, description, reporter, state, make(map[string]string, 0)}
		}
		if attname.Valid && attvalue.Valid {
			retVal.Attributes[attname.String] = attvalue.String
		}
	}

	logManager.LogPrintf("got incident: %v\n", retVal)
	return retVal, true
}

func (manager MySQLManager) GetIncidents() ([]Incident, bool) {
	incidents := make(map[int64]Incident, 0)
	var (
		id           int64
		incidenttype string
		description  string
		reporter     string
		state        string
		attname      sql.NullString
		attvalue     sql.NullString
	)

	if manager.Connection == nil {
		logManager.LogFatalln("Connection is nil")
	}

	rows, err := manager.Connection.Query("SELECT Id, Type, Description, Reporter, State, AttributeName, AttributeValue " +
		"FROM incidents LEFT JOIN incidentattributes " +
		"ON IncidentId = Id")

	if err != nil {
		logManager.LogPrintf("Error occurred when preparing get %v\n", err)
		return nil, false
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &incidenttype, &description, &reporter, &state, &attname, &attvalue)
		if err != nil {
			logManager.LogPrintln(err)
		}

		val, found := incidents[id]
		if !found {
			val = Incident{incidenttype, id, description, reporter, state, make(map[string]string, 0)}
		}

		if attname.Valid && attvalue.Valid {
			val.Attributes[attname.String] = attvalue.String
		}

		incidents[id] = val
	}

	retVal := make([]Incident, len(incidents))
	iter := 0
	for _, value := range incidents {
		retVal[iter] = value
		iter++
	}
	logManager.LogPrintf("got incidents: %v\n", retVal)
	return retVal, true
}

func (manager MySQLManager) UpdateIncident(id int, incident IncidentUpdate) bool {
	return true
}

func (manager MySQLManager) AddAttachment(incidentId int, attachment Attachment) bool {
	return true
}

func (manager MySQLManager) GetAttachments(incidentId int) ([]Attachment, bool) {
	return make([]Attachment, 0), true
}

func (manager MySQLManager) RemoveAttachment(incidentId int, fileName string) bool {
	return true
}

// CleanUp will do any required cleanup actions on the incident manager.
func (manager MySQLManager) CleanUp() {
	logManager.LogPrintln("Closing database connection")
	if manager.Connection != nil {
		manager.Connection.Close()
	}
}
