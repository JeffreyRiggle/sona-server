package main

type UserManager interface {
	AddUser(user *AddUser) (bool, User)
	GetUser(userId int64) (User, bool)
	UpdateUser(userId int64, user *User) bool
	RemoveUser(userId int64) bool
	SetUserPassword(user User, password string)
	SetPermissions(userId int64, permissions []string) bool
	AuthenticateUser(user User, password string) (bool, TokenResponse)
	ValidateUser(token string) bool
}
