package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type MySQLUserManager struct {
	Connection         *sql.DB
	DefaultPermissions []string
}

func (manager MySQLUserManager) Initialize() {
	if !manager.hasTable("Users") {
		logManager.LogPrintln("Unable to find incident table creating now")
		manager.createUserTable()
	}
}

func (manager MySQLUserManager) hasTable(tableName string) bool {
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

func (manager MySQLUserManager) createUserTable() {
	stmt, err := manager.Connection.Prepare("CREATE TABLE Users (" +
		"Id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY, " +
		"EmailAddress VARCHAR(1048), " +
		"UserName VARCHAR(255), " +
		"FirstName VARCHAR(255), " +
		"Lastname VARCHAR(255), " +
		"Gender VARCHAR(255), " +
		"Password VARCHAR(1048), " +
		"Permissions VARCHAR(1048))")

	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec()
	if err != nil {
		panic(err)
	}

	logManager.LogPrintln("Created Incident Table: %v\n", res)
}

func (manager MySQLUserManager) AddUser(user *AddUser) (bool, User) {
	stmt, err := manager.Connection.Prepare("INSERT INTO Users (EmailAddress, UserName, FirstName, LastName, Gender, Password, Permissions) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		logManager.LogPrintf("Error occurred when preparing add %v", err)
		return false, User{}
	}

	permissions := make([]string, len(manager.DefaultPermissions))

	res, err := stmt.Exec(user.EmailAddress, user.UserName, user.FirstName, user.LastName, user.Gender, user.Password, strings.Join(permissions, ","))

	if err != nil {
		logManager.LogPrintf("Error occurred when executing add %v", err)
		return false, User{}
	}

	id, err := res.LastInsertId()

	if err != nil {
		logManager.LogPrintf("Unable to get last inserted id %v\n", err)
		return false, User{}
	}

	usr := User{
		Id:           id,
		UserName:     user.UserName,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		EmailAddress: user.EmailAddress,
		Gender:       user.Gender,
		Permissions:  permissions,
	}

	logManager.LogPrintf("Created new user: %v\n", usr)
	return true, usr
}

func (manager MySQLUserManager) GetUser(userId int64) (User, bool) {
	retVal := User{}
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
		"FROM Users "+
		"WHERE Id = ?", userId)

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
			retVal = User{}
		}
	}

	logManager.LogPrintf("got user: %v\n", retVal)
	return retVal, true
}

func (manager MySQLUserManager) UpdateUser(userId int64, user *User) bool {
	// TODO
	return true
}

func (manager MySQLUserManager) RemoveUser(userId int64) bool {
	// TODO
	return true
}

func (manager MySQLUserManager) SetUserPassword(user User, password string) {
	// TODO
}

func (manager MySQLUserManager) SetPermissions(userId int64, permissions []string) bool {
	// TODO
	return true
}

func (manager MySQLUserManager) AuthenticateUser(user User, password string) (bool, TokenResponse) {
	// TODO
	return true, TokenResponse{}
}

func (manager MySQLUserManager) ValidateUser(token string) bool {
	// TODO
	return true
}

// CleanUp will do any required cleanup actions on the incident manager.
func (manager MySQLUserManager) CleanUp() {
	logManager.LogPrintln("Closing database connection")
	if manager.Connection != nil {
		manager.Connection.Close()
	}
}
