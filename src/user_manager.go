package main

type UserManager interface {
	AddUser(user *AddUser) (bool, User)
	GetUser(userId int) (User, bool)
	UpdateUser(userId int, user *User) bool
	RemoveUser(userId int) bool
	SetUserPassword(user User, password string)
	AuthenticateUser(user User, password string) bool
}
