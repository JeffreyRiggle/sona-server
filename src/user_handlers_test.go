package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"github.com/gorilla/mux"
)

var usrRouter *mux.Router

func userTestSetup() {
	if router == nil {
		usrRouter = NewRouter()
		http.Handle("/", usrRouter)
	}

	userManager = RuntimeUserManager{make(map[int]*User), make(map[int]string), make(map[int][]string), make([]string, 0)}
	hookManager = HookManager{make([]WebHook, 0), make([]WebHook, 0), make([]WebHook, 0), make([]WebHook, 0), make([]WebHook, 0)}
}

func TestCreateUser(t *testing.T) {
	userTestSetup()
	usr := AddUser{
		EmailAddress: "a@b.c",
		FirstName:    "Foo",
		LastName:     "User",
		UserName:     "FooUser",
		Password:     "1234",
	}

	body, _ := json.Marshal(usr)

	r, _ := http.NewRequest("POST", "/sona/v1/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	usrRouter.ServeHTTP(w, r)

	if w.Result().StatusCode != 201 {
		t.Errorf("Expected 201 status code got %v", w.Result())
	}

	var retVal User
	err := json.Unmarshal(w.Body.Bytes(), &retVal)
	if err != nil {
		t.Errorf("Failed to convert response %v error %v", w.Body, err)
	}

	if retVal.Id != 0 {
		t.Errorf("Expected user 0 got %v", retVal.Id)
	}

	if retVal.EmailAddress != "a@b.c" {
		t.Errorf("Expected email address a@b.c got %v", retVal.EmailAddress)
	}

	if retVal.FirstName != "Foo" {
		t.Errorf("Expected first name Foo got %v", retVal.FirstName)
	}

	if retVal.LastName != "User" {
		t.Errorf("Expected last name User got %v", retVal.LastName)
	}

	if retVal.UserName != "FooUser" {
		t.Errorf("Expected user name FooUser got %v", retVal.UserName)
	}
}

func TestCreateUserWithBadData(t *testing.T) {
	userTestSetup()
	usr := AddUser{}

	body, _ := json.Marshal(usr)

	r, _ := http.NewRequest("POST", "/sona/v1/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	usrRouter.ServeHTTP(w, r)

	if w.Result().StatusCode != 400 {
		t.Errorf("Expected 400 status code got %v", w.Result())
	}
}
