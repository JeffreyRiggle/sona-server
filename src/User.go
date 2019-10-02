package main

type AddUser struct {
	EmailAddress string `json:"emailAddress"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Gender       string `json:"gender"`
	Password     string `json:"password"`
}

type User struct {
	EmailAddress string `json:"emailAddress"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Gender       string `json:"gender"`
	Id           int    `json:"id"`
}

type UserPassword struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type PasswordChangeRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (user User) SetPassword(password string) {
	userManager.SetUserPassword(user, password)
}

func (user User) Authenticate(password string) (bool, TokenResponse) {
	return userManager.AuthenticateUser(user, password)
}