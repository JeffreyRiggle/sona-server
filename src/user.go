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
	Id           int64    `json:"id"`
	Permissions  []string `json:"permissions"`
}

type UserPassword struct {
	EmailAddress string `json:"emailAddress"`
	Password     string `json:"password"`
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

func updateUser(original *User, updated User) bool {
	changed := false
	if len(updated.UserName) > 0 {
		original.UserName = updated.UserName
		changed = true
	}

	if len(updated.FirstName) > 0 {
		original.FirstName = updated.FirstName
		changed = true
	}

	if len(updated.LastName) > 0 {
		original.LastName = updated.LastName
		changed = true
	}

	if len(updated.Gender) > 0 {
		original.Gender = updated.Gender
		changed = true
	}

	if updated.Permissions != nil {
		original.Permissions = updated.Permissions
		changed = true
	}

	return changed
}

func getUserPropertyValue(key string, user User) string {
	if strings.EqualFold(key, "id") {
		return strconv.FormatInt(user.Id, 10)
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
