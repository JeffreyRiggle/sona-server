package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var userManager UserManager

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Got Create User request")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	addUser, pass := convertAddUser(r.Body)
	if !pass {
		logManager.LogPrintf("Bad request for create user %v", addUser)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	passed, user := userManager.AddUser(&addUser)
	if !passed {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logManager.LogPrintf("Created user %v\n", user.Id)
	go hookManager.CallAddedUserHooks(user)

	data, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func convertAddUser(body io.ReadCloser) (AddUser, bool) {
	decoder := json.NewDecoder(body)

	var user AddUser
	err := decoder.Decode(&user)

	if err != nil {
		logManager.LogPrintf("Got error when attempting to decode body %v", err)
		return user, false
	}

	if len(user.EmailAddress) == 0 || len(user.UserName) == 0 || len(user.Password) == 0 {
		return user, false
	}

	return user, true
}

func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrint("Got User Update")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("X-Sona-Token")
	if !userManager.ValidateUser(token) {
		logManager.LogPrintf("Invalid Token %v used", token)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["userId"])

	if err != nil {
		logManager.LogPrintf("Error converting userId %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !HasPermission(token, availablePermissions.modifyUser) && GetTokenUser(token) != userId {
		logManager.LogPrintf("Token does not allow for modify user", token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	update, passed := convertUpdateUser(r.Body)

	if !passed {
		logManager.LogPrintf("Invalid update for %v\n", userId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userManager.UpdateUser(userId, &update) {
		go hookManager.CallUpdatedUserHooks(update)
		w.WriteHeader(http.StatusOK)
		return
	}

	logManager.LogPrintf("User %v not found\n", userId)
	w.WriteHeader(http.StatusNotFound)
}

func convertUpdateUser(body io.ReadCloser) (User, bool) {
	decoder := json.NewDecoder(body)

	var update User
	err := decoder.Decode(&update)

	if err != nil {
		return update, false
	}

	return update, true
}

func HandleSetPermissions(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrint("Got User Update")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("X-Sona-Token")
	if !userManager.ValidateUser(token) {
		logManager.LogPrintf("Invalid Token %v used", token)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !HasPermission(token, availablePermissions.master) {
		logManager.LogPrintf("Token does not allow for modify user", token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["userId"])

	if err != nil {
		logManager.LogPrintf("Error converting userId %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var permissions []string
	err2 := decoder.Decode(&permissions)

	if err2 != nil {
		logManager.LogPrintf("Error decoding new permissions %v", err2)
		panic(err2)
	}

	logManager.LogPrintf("Attempting to set permissions to %v", permissions)

	if userManager.SetPermissions(userId, permissions) {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Got user state request")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	token := r.Header.Get("X-Sona-Token")
	if !userManager.ValidateUser(token) {
		logManager.LogPrintf("Invalid Token %v used", token)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["userId"])
	if err != nil {
		logManager.LogPrintf("Error converting userId %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !HasPermission(token, availablePermissions.viewUser) && GetTokenUser(token) != userId {
		logManager.LogPrintf("Token does not allow for view user", token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if val, ok := userManager.GetUser(userId); ok {
		logManager.LogPrintf("Got State request for %v.", userId)
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(val); err != nil {
			panic(err)
		}

		return
	}

	logManager.LogPrintf("User %v not found\n", userId)
	w.WriteHeader(http.StatusNotFound)
}

func HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Got user delete request")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	token := r.Header.Get("X-Sona-Token")
	if !userManager.ValidateUser(token) {
		logManager.LogPrintf("Invalid Token %v used", token)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if !HasPermission(token, availablePermissions.deleteUser) {
		logManager.LogPrintf("Token does not allow for modify user", token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["userId"])

	if err != nil {
		logManager.LogPrintf("Error converting userId %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if userManager.RemoveUser(userId) {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}

func HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrintln("Got user password change request.")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	token := r.Header.Get("X-Sona-Token")
	if !userManager.ValidateUser(token) {
		logManager.LogPrintf("Invalid Token %v used", token)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)

	userId, err := strconv.Atoi(vars["userId"])

	if err != nil {
		logManager.LogPrintf("Error converting userId %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !HasPermission(token, availablePermissions.master) && GetTokenUser(token) != userId {
		logManager.LogPrintf("Token does not allow for modify user", token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var req PasswordChangeRequest
	err2 := decoder.Decode(&req)

	if err2 != nil {
		logManager.LogPrint("Invalid password change request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, found := userManager.GetUser(userId)

	if !found {
		logManager.LogPrintf("Error finding userId %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	auth, _ := user.Authenticate(req.OldPassword)

	if !auth {
		logManager.LogPrintf("Failed to authenticate %v", userId)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	user.SetPassword(req.NewPassword)
	w.WriteHeader(http.StatusOK)
}

func HandleAuthentication(w http.ResponseWriter, r *http.Request) {
	logManager.LogPrint("Got User Authentication")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)

	var req UserPassword
	err := decoder.Decode(&req)

	if err != nil {
		logManager.LogPrint("Invalid password request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, found := userManager.GetUser(req.Id)

	if !found {
		logManager.LogPrintf("Unable to find user %v", req.Id)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	auth, token := user.Authenticate(req.Password)

	if !auth {
		logManager.LogPrintf("Failed to authenticate %v", req.Id)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(token); err != nil {
		panic(err)
	}
}
