package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route defines a url endpoint and its handler.
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// Routes defines a collection of Route.
type Routes []Route

// NewRouter creates a new mux.Router with the defined Routes.
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.Handler)
	}

	return router
}

// TODO create SQL user manager
// TODO create dynamodb user manager
// TODO create datastore user manager
var routes = Routes{
	Route{
		"Create",
		"POST",
		"/sona/v1/incidents",
		HandleCreateIncident,
	},
	Route{
		"Update",
		"PUT",
		"/sona/v1/incidents/{incidentId}",
		HandleIncidentUpdate,
	},
	Route{
		"GetAttachments",
		"GET",
		"/sona/v1/incidents/{incidentId}/attachments",
		HandleGetAttachments,
	},
	Route{
		"UploadAttachment",
		"POST",
		"/sona/v1/incidents/{incidentId}/attachment",
		HandleUploadAttachment,
	},
	Route{
		"DownloadAttachment",
		"GET",
		"/sona/v1/incidents/{incidentId}/attachment/{attachmentId}",
		HandleDownloadAttachment,
	},
	Route{
		"RemoveAttachment",
		"DELETE",
		"/sona/v1/incidents/{incidentId}/attachment/{attachmentId}",
		HandleRemoveAttachment,
	},
	Route{
		"GetIncidents",
		"GET",
		"/sona/v1/incidents",
		HandleGetIncidents,
	},
	Route{
		"GetIncident",
		"GET",
		"/sona/v1/incidents/{incidentId}",
		HandleGetIncident,
	},
	Route{
		"CreateUser",
		"POST",
		"/sona/v1/users",
		HandleCreateUser,
	},
	Route{
		"GetUser",
		"GET",
		"/sona/v1/users/{userId}",
		HandleGetUser,
	},
	Route{
		"UpdateUser",
		"PUT",
		"/sona/v1/users/{userId}",
		HandleUpdateUser,
	},
	Route{
		"DeleteUser",
		"DELETE",
		"/sona/v1/users/{userId}",
		HandleDeleteUser,
	},
	Route{
		"Authentication",
		"PUT",
		"/sona/v1/users/{userId}/authentication",
		HandleChangePassword,
	},
	Route{
		"UserPermissions",
		"PUT",
		"/sona/v1/users/{userId}/permissions",
		HandleSetPermissions,
	},
	Route{
		"Authenticate",
		"POST",
		"/sona/v1/authenticate",
		HandleAuthentication,
	},
}
