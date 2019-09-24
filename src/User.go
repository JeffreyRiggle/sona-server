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

func (user User) setPassword(password string) {
	userManager.SetUserPassword(user.Id, password)
}

func (user User) authenticate(password string) bool {
	return userManager.AuthenticateUser(user.Id, password)
}