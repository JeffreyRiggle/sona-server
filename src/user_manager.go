package main

type UserManager interface {
	AddUser(user *AddUser) (bool, User)
	GetUser(userId int) (User, bool)
	UpdateUser(userId int, user *User) bool
	RemoveUser(userId int) bool
	SetUserPassword(user User, password string)
	SetPermissions(userId int, permissions []string) bool
	AuthenticateUser(user User, password string) (bool, TokenResponse)
	ValidateUser(token string) bool
}
