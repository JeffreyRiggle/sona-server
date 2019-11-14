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
		logManager.LogPrintln("Unable to find user table creating now")
		manager.createUserTable()
	}

	if !manager.hasTable("Tokens") {
		logManager.LogPrintln("Unable to find tokens table creating now")
		manager.createTokenTable()
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
		"Password BLOB(65535), " + // TODO figure this out.
		"Permissions VARCHAR(1048))")

	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec()
	if err != nil {
		panic(err)
	}

	logManager.LogPrintln("Created User Table: %v\n", res)
}

func (manager MySQLUserManager) createTokenTable() {
	stmt, err := manager.Connection.Prepare("CREATE TABLE Tokens (" +
		"Id INT UNSIGNED NOT NULL PRIMARY KEY, " +
		"Tokens VARCHAR(1048))")

	if err != nil {
		panic(err)
	}

	res, err := stmt.Exec()
	if err != nil {
		panic(err)
	}

	logManager.LogPrintln("Created Token Table: %v\n", res)
}

func (manager MySQLUserManager) AddUser(user *AddUser) (bool, User) {
	stmt, err := manager.Connection.Prepare("INSERT INTO Users (EmailAddress, UserName, FirstName, LastName, Gender, Permissions) " +
		"VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		logManager.LogPrintf("Error occurred when preparing add %v", err)
		return false, User{}
	}

	permissions := make([]string, len(manager.DefaultPermissions))

	res, err := stmt.Exec(user.EmailAddress, user.UserName, user.FirstName, user.LastName, user.Gender, strings.Join(permissions, ","))

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

	manager.SetUserPassword(usr, createPasswordHash(usr, user.Password))

	logManager.LogPrintf("Created new user: %v\n", usr)
	return true, usr
}

func (manager MySQLUserManager) GetUser(userId int64) (User, bool) {
	retVal := User{
		Id: -1,
	}
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

		if retVal.Id == -1 {
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
	}

	logManager.LogPrintf("got user: %v\n", retVal)
	return retVal, retVal.Id != -1
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
	_, pass := manager.GetUser(userId)

	if !pass {
		return false
	}

	stmt, err := manager.Connection.Prepare("UPDATE Users SET Permissions = ? WHERE Id = ?")
	if err != nil {
		logManager.LogPrintf("Error occurred when preparing update user permissions %v", err)
		return false
	}

	_, err = stmt.Exec(strings.Join(permissions, ","), userId)

	if err != nil {
		logManager.LogPrintf("Error occurred when executing update user permissions %v", err)
		return false
	}

	return true
}

func (manager MySQLUserManager) AuthenticateUser(user User, password string) (bool, TokenResponse) {
	var storedPassword string

	rows, err := manager.Connection.Query("SELECT Password "+
		"FROM Users "+
		"WHERE Id = ?", user.Id)

	if err != nil {
		logManager.LogPrintf("Error occurred when preparing get password %v\n", err)
		return false, TokenResponse{}
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&storedPassword)
		if err != nil {
			logManager.LogPrintln(err)
		}
	}

	auth := createPasswordHash(user, password) == storedPassword

	if !auth {
		return auth, TokenResponse{}
	}

	token := GenerateToken(user)
	tokens, found := manager.reconcileUserTokens(user.Id, token.Token)

	if found {
		manager.updateTokens(user.Id, tokens)
	} else {
		manager.setTokens(user.Id, tokens)
	}

	return auth, token
}

func (manager MySQLUserManager) reconcileUserTokens(userId int64, newToken string) (string, bool) {
	tokens, found := manager.getTokens(userId)
	reconciledTokens := make([]string, 0)
	reconciledTokens = append(reconciledTokens, newToken)

	for _, token := range tokens {
		if !TokenExpired(token) {
			reconciledTokens = append(reconciledTokens, token)
		}
	}

	return strings.Join(reconciledTokens, ","), found
}

func (manager MySQLUserManager) getTokens(userId int64) ([]string, bool) {
	var (
		id          int64
		storedToken string
	)

	rows, err := manager.Connection.Query("SELECT Id, Tokens "+
		"FROM Tokens "+
		"WHERE Id = ?", userId)

	if err != nil {
		logManager.LogPrintf("Error occurred when preparing get tokens %v\n", err)
		return make([]string, 0), false
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &storedToken)
		if err != nil {
			logManager.LogPrintln(err)
		}
	}

	return strings.Split(storedToken, ","), id != 0
}

func (manager MySQLUserManager) updateTokens(userId int64, tokens string) {
	stmt, err := manager.Connection.Prepare("UPDATE Tokens SET Tokens = ? WHERE Id = ?")
	if err != nil {
		logManager.LogPrintf("Error occurred when preparing update user tokens %v", err)
		return
	}

	_, err = stmt.Exec(tokens, userId)

	if err != nil {
		logManager.LogPrintf("Error occurred when executing update user tokens %v", err)
	}
}

func (manager MySQLUserManager) setTokens(userId int64, tokens string) {
	stmt, err := manager.Connection.Prepare("INSERT INTO Tokens (Id, Tokens) VALUES (?, ?)")
	if err != nil {
		logManager.LogPrintf("Error occurred when preparing add user tokens %v", err)
		return
	}

	_, err = stmt.Exec(userId, tokens)

	if err != nil {
		logManager.LogPrintf("Error occurred when executing add user tokens %v", err)
	}
}

func (manager MySQLUserManager) ValidateUser(token string) bool {
	userId := GetTokenUser(token)
	tokens, _ := manager.getTokens(userId)
	found := -1

	for i, v := range tokens {
		if v == token {
			found = i
		}
	}

	if found == -1 {
		logManager.LogPrintf("Token not found for user %v", userId)
		return false
	}

	expired := TokenExpired(token)
	logManager.LogPrintf("Token expired %v", expired)
	go manager.pruneTokens(userId, tokens)

	return !expired
}

func (manager MySQLUserManager) pruneTokens(userId int64, tokens []string) {
	prunedTokens := make([]string, 0)

	for _, v := range tokens {
		if !TokenExpired(v) {
			prunedTokens = append(prunedTokens, v)
		}
	}

	if len(tokens) != len(prunedTokens) {
		manager.setTokens(userId, strings.Join(prunedTokens, ","))
	}
}

// CleanUp will do any required cleanup actions on the user manager.
func (manager MySQLUserManager) CleanUp() {
	logManager.LogPrintln("Closing database connection")
	if manager.Connection != nil {
		manager.Connection.Close()
	}
}
