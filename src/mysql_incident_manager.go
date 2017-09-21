package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLManager struct {
	UserName   string
	Password   string
	Host       string
	Port       string
	DBName     string
	Connection *sql.DB
}

func (manager MySQLManager) Initialize() {
	conn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", manager.UserName, manager.Password, manager.Host, manager.Port, manager.DBName)
	db, err := sql.Open("mysql", conn)

	if err != nil {
		panic(err)
	}

	manager.Connection = db

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
		"VALUES (?, ?, ?, ?)")
	if err != nil {
		return false
	}

	res, err := stmt.Exec(incident.Type, incident.Description, incident.Reporter, incident.State)

	if err != nil {
		return false
	}

	logManager.LogPrintf("Created new incident: %v\n", res)
	return true
}

func (manager MySQLManager) GetIncident(incidentId int) (Incident, bool) {
	return Incident{}, true
}

func (manager MySQLManager) GetIncidents() ([]Incident, bool) {
	return make([]Incident, 0), true
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
	if manager.Connection != nil {
		manager.Connection.Close()
	}
}
