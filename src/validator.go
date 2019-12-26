package main

import "net/http"

func validateRequest(w http.ResponseWriter, r *http.Request, permission string) bool {
	token := r.Header.Get("X-Sona-Token")
	param := r.URL.Query()["token"]

	if len(token) == 0 && param != nil {
		logManager.LogPrintf("Using query parameter over header")
		token = param[0]
	}

	if !userManager.ValidateUser(token) {
		logManager.LogPrintf("Invalid Token %v used", token)
		w.WriteHeader(http.StatusForbidden)
		return false
	}

	if !HasPermission(token, permission) {
		logManager.LogPrintf("Token %v does not allow for %v", token, permission)
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}
