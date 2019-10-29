package main

import (
	"strconv"
	"strings"
)

type AddUser struct {
	EmailAddress string `json:"emailAddress"`
	UserName     string `json:"userName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Gender       string `json:"gender"`
	Password     string `json:"password"`
}

type User struct {
	EmailAddress string   `json:"emailAddress"`
	UserName     string   `json:"userName"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	Gender       string   `json:"gender"`
	Id           int      `json:"id"`
	Permissions  []string `json:"permissions"`
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

func getUserPropertyValue(key string, user User) string {
	if strings.EqualFold(key, "id") {
		return strconv.Itoa(user.Id)
	}
	if strings.EqualFold(key, "username") {
		return user.UserName
	}
	if strings.EqualFold(key, "firstname") {
		return user.FirstName
	}
	if strings.EqualFold(key, "lastname") {
		return user.LastName
	}
	if strings.EqualFold(key, "gender") {
		return user.Gender
	}
	if strings.EqualFold(key, "emailaddress") {
		return user.EmailAddress
	}

	return ""
}
