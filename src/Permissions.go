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

// TODO start enforcing permissions
// TODO create default administrator user that has master
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
