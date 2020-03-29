package main

import (
	b64 "encoding/base64"
	"strconv"
	"strings"
	"time"

	guuid "github.com/google/uuid"
)

type TokenResponse struct {
	Token  string `json:"token"`
	UserId int64  `json:"userId"`
}

func GenerateToken(user User) TokenResponse {
	now := time.Now()
	timeout := now.Add(time.Hour * time.Duration(3))

	id := guuid.New().String()
	logManager.LogPrintf("Token will expire at %v\n", timeout)

	permissions := strings.Join(user.Permissions, ",")
	val := strconv.FormatInt(user.Id, 10) + ":" + id + ":" + strconv.FormatInt(timeout.UnixNano(), 10) + ":" + permissions
	return TokenResponse{b64.StdEncoding.EncodeToString([]byte(val)), user.Id}
}

func GetTokenUser(token string) int64 {
	decoded, _ := b64.StdEncoding.DecodeString(token)
	vals := strings.Split(string(decoded), ":")

	if len(vals) < 4 {
		return -1
	}

	retVal, _ := strconv.ParseInt(vals[0], 10, 64)

	return retVal
}

func TokenExpired(token string) bool {
	decoded, _ := b64.StdEncoding.DecodeString(token)
	vals := strings.Split(string(decoded), ":")

	if len(vals) < 4 {
		return false
	}

	timestamp, err := strconv.ParseInt(vals[2], 10, 64)

	if err != nil {
		logManager.LogFatal(err)
		return false
	}

	return timestamp < time.Now().UnixNano()
}

func HasPermission(token string, permission string) bool {
	decoded, _ := b64.StdEncoding.DecodeString(token)
	vals := strings.Split(string(decoded), ":")

	if len(vals) < 4 {
		return false
	}

	permissions := strings.Split(vals[3], ",")

	for _, p := range permissions {
		if p == permission || p == availablePermissions.master {
			return true
		}
	}

	return false
}
