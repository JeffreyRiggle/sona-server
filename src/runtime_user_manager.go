package main

import (
	"crypto/sha256"
	"io"
	"strings"
)

type RuntimeUserManager struct {
	Users              map[int64]*User
	Passwords          map[int64]string
	Tokens             map[int64][]string
	DefaultPermissions []string
}

func (manager RuntimeUserManager) AddUser(user *AddUser) (bool, User) {
	var cuser, id = manager.convertAddUser(user)
	manager.Users[id] = cuser

	manager.SetUserPassword(*cuser, user.Password)
	return true, *cuser
}

func (manager RuntimeUserManager) convertAddUser(user *AddUser) (*User, int64) {
	var retVal User
	var id = int64(len(manager.Users))

	retVal.Permissions = make([]string, len(manager.DefaultPermissions))
	copy(retVal.Permissions, manager.DefaultPermissions)
	retVal.Id = id
	retVal.EmailAddress = user.EmailAddress
	retVal.UserName = user.UserName
	retVal.FirstName = user.FirstName

	if len(user.LastName) != 0 {
		retVal.LastName = user.LastName
	}

	if len(user.Gender) != 0 {
		retVal.Gender = user.Gender
	}

	return &retVal, id
}

func (manager RuntimeUserManager) GetUser(userId int64) (User, bool) {
	if val, ok := manager.Users[userId]; ok {
		return *val, true
	}

	return User{}, false
}

func (manager RuntimeUserManager) GetUserByEmail(emailAddress string) (User, bool) {
	for _, u := range manager.Users {
		if u.EmailAddress == emailAddress {
			return *u, true
		}
	}

	return User{}, false
}

func (manager RuntimeUserManager) UpdateUser(userId int64, user *User) bool {
	originalUser := manager.Users[userId]
	updateUser(originalUser, *user)
	manager.Users[userId] = originalUser
	return true
}

func (manager RuntimeUserManager) RemoveUser(userId int64) bool {
	delete(manager.Users, userId)
	return true
}

func (manager RuntimeUserManager) SetUserPassword(user User, password string) {
	manager.Passwords[user.Id] = createPasswordHash(user, password)
}

func (manager RuntimeUserManager) AuthenticateUser(user User, password string) (bool, TokenResponse) {
	auth := createPasswordHash(user, password) == manager.Passwords[user.Id]

	if !auth {
		return auth, TokenResponse{"", -1}
	}

	token := GenerateToken(user)

	if manager.Tokens[user.Id] == nil {
		manager.Tokens[user.Id] = make([]string, 0, 0)
	}

	manager.Tokens[user.Id] = append(manager.Tokens[user.Id], token.Token)
	return auth, token
}

func (manager RuntimeUserManager) ValidateUser(token string) bool {
	userId := GetTokenUser(token)

	found := -1

	for i, v := range manager.Tokens[userId] {
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

	if expired {
		manager.Tokens[userId] = append(manager.Tokens[userId][:found], manager.Tokens[userId][found+1:]...)
	}

	return !expired
}

func (manager RuntimeUserManager) SetPermissions(userId int64, permissions []string) bool {
	user := manager.Users[userId]
	user.Permissions = permissions
	manager.Users[userId] = user
	return true
}

func createPasswordHash(user User, password string) string {
	hash := sha256.New()

	pw := strings.NewReader(password + user.EmailAddress)

	if _, err := io.Copy(hash, pw); err != nil {
		return ""
	}

	return string(hash.Sum(nil))
}
