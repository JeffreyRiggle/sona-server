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
		username     string
		firstname    string
		lastname     string
		emailaddress string
		gender       string
		permissions  string
	)

	rows, err := manager.Connection.Query("SELECT Id, UserName, FirstName, LastName, EmailAddress, Gender, Permissions "+
		"FROM Users "+
		"WHERE Id = ?", userId)

	if err != nil {
		logManager.LogPrintf("Error occurred when preparing get %v\n", err)
		return retVal, false
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &username, &firstname, &lastname, &emailaddress, &gender, &permissions)
		if err != nil {
			logManager.LogPrintln(err)
		}

		retVal = User{
			Id:           id,
			UserName:     username,
			FirstName:    firstname,
			LastName:     lastname,
			EmailAddress: emailaddress,
			Gender:       gender,
			Permissions:  strings.Split(permissions, ","),
		}
	}

	logManager.LogPrintf("got user: %v\n", retVal)
	return retVal, true
}

func (manager MySQLUserManager) UpdateUser(userId int64, user *User) bool {
	usr, pass := manager.GetUser(userId)

	if !pass {
		return false
	}

	if !updateUser(&usr, *user) {
		return true
	}

	stmt, err := manager.Connection.Prepare("UPDATE Users SET UserName = ? , FirstName = ?, LastName = ?, Gender = ? WHERE Id = ?")
	if err != nil {
		logManager.LogPrintf("Error occurred when preparing update user %v", err)
		return false
	}

	_, err = stmt.Exec(usr.UserName, usr.FirstName, usr.LastName, usr.Gender, strings.Join(usr.Permissions, ","), userId)

	if err != nil {
		logManager.LogPrintf("Error occurred when executing update user %v", err)
		return false
	}

	return true
}

func (manager MySQLUserManager) RemoveUser(userId int64) bool {
	stmt, err := manager.Connection.Prepare("DELETE FROM Users WHERE Id = ?")
	if err != nil {
		logManager.LogPrintf("Error occurred when preparing remove user %v", err)
		return false
	}

	_, err = stmt.Exec(userId)

	if err != nil {
		logManager.LogPrintf("Error occurred when executing remove user %v", err)
		return false
	}

	return true
}

func (manager MySQLUserManager) SetUserPassword(user User, password string) {
	stmt, err := manager.Connection.Prepare("UPDATE Users SET Password = ? WHERE Id = ?")
	if err != nil {
		logManager.LogPrintf("Error occurred when preparing update user password %v", err)
		return
	}

	_, err = stmt.Exec(password, user.Id)

	if err != nil {
		logManager.LogPrintf("Error occurred when executing update user password %v", err)
		return
	}
}

func (manager MySQLUserManager) SetPermissions(userId int64, permissions []string) bool {
	usr, pass := manager.GetUser(userId)

	if !pass {
		return false
	}

	stmt, err := manager.Connection.Prepare("UPDATE Users SET Permissions = ? WHERE Id = ?")
	if err != nil {
		logManager.LogPrintf("Error occurred when preparing update user permissions %v", err)
		return false
	}

	_, err = stmt.Exec(strings.Join(usr.Permissions, ","), userId)

	if err != nil {
		logManager.LogPrintf("Error occurred when executing update user permissions %v", err)
		return false
	}

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

// CleanUp will do any required cleanup actions on the user manager.
func (manager MySQLUserManager) CleanUp() {
	logManager.LogPrintln("Closing database connection")
	if manager.Connection != nil {
		manager.Connection.Close()
	}
}
