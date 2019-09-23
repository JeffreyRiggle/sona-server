package main

type UserManager interface {
	AddUser(user *User) bool
	GetUser(userId int) (User, bool)
	UpdateUser(userId int, user *User) bool
	RemoveUser(userId int) bool
	SetUserPassword(userId int, password string)
	AuthenticateUser(userId int, password string) bool
}
