package main

import (
	b64 "encoding/base64"
	"strconv"
	"strings"
	"time"
)

type TokenResponse struct {
	Token string `json:"token"`
}

func GenerateToken(user User) TokenResponse {
	now := time.Now()
	timeout := now.Add(time.Hour * time.Duration(3))

	logManager.LogPrintf("Token will expire at %v\n", timeout)

	val := strconv.Itoa(user.Id) + ":" + strconv.FormatInt(timeout.UnixNano(), 10)
	return TokenResponse{b64.StdEncoding.EncodeToString([]byte(val))}
}

func TokenExpired(token string) bool {
	decoded, _ := b64.StdEncoding.DecodeString(token)
	vals := strings.Split(string(decoded), ":")

	if len(vals) < 2 {
		return false
	}

	timestamp, err := strconv.ParseInt(vals[1], 10, 64)

	if err != nil {
		return false
	}
	
	return timestamp > time.Now().UnixNano()
}