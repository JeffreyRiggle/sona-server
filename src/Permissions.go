package main

type AvailablePermissions struct {
	viewIncident   string
	modifyIncident string
	deleteIncident string
	createIncident string
	viewUser       string
	modifyUser     string
	deleteUser     string
	master         string
}

var availablePermissions = AvailablePermissions{
	"incident-view",
	"incident-modify",
	"incident-delete",
	"incident-create",
	"user-view",
	"user-modify",
	"user-delete",
	"*",
}
