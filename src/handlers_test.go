package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"encoding/json"

	"github.com/gorilla/mux"
)

var router *mux.Router
var user1 User

// FakeFileManager a manager to fake file io
type FakeFileManager struct {
}

func (manager FakeFileManager) SaveFile(incident string, fileName string, file multipart.File) (string, bool) {
	return fileName, true
}

func (manager FakeFileManager) LoadFile(incident string, fileName string) (io.ReadSeeker, os.FileInfo, bool, func()) {
	return nil, nil, true, nil
}

func (manager FakeFileManager) DeleteFile(incident string, fileName string) bool {
	return true
}

func setup() {
	if router == nil {
		router = NewRouter()
		http.Handle("/", router)
	}

	incidentManager = RuntimeIncidentManager{make(map[int64]*Incident), make(map[int][]Attachment)}
	userManager = RuntimeUserManager{make(map[int]*User), make(map[int]string), make(map[int][]string)}
	hookManager = HookManager{make([]WebHook, 0), make([]WebHook, 0), make([]WebHook, 0), make([]WebHook, 0), make([]WebHook, 0)}
	fileManager = FakeFileManager{}

	addUser1 := AddUser{
		EmailAddress: "a@b.c",
		FirstName:    "Foo",
		LastName:     "User",
		UserName:     "FooUser",
		Password:     "1234",
	}

	_, user1 = userManager.AddUser(&addUser1)
}

func TestCreateIncident(t *testing.T) {
	setup()
	inc := Incident{}
	inc.Reporter = "Tester"
	inc.Description = "Some Test"
	body, _ := json.Marshal(inc)

	r, _ := http.NewRequest("POST", "/sona/v1/create", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 201 {
		t.Errorf("Expected 201 status code got %v", w.Result())
	}

	var retVal Incident
	err := json.Unmarshal(w.Body.Bytes(), &retVal)
	if err != nil {
		t.Errorf("Failed to convert response %v error %v", w.Body, err)
	}

	if retVal.Id != 0 {
		t.Errorf("Expected incident 0 got %v", retVal.Id)
	}

	if retVal.Reporter != "Tester" {
		t.Errorf("Expected reporter Tester got %v", retVal.Reporter)
	}

	if retVal.Description != "Some Test" {
		t.Errorf("Expected description Some Test got %v", retVal.Description)
	}

	if retVal.State != "open" {
		t.Errorf("Expected state open got %v", retVal.State)
	}
}

func TestCreateIncidentWithInvalidRequest(t *testing.T) {
	setup()
	att := Attachment{}
	body, _ := json.Marshal(att)

	r, _ := http.NewRequest("POST", "/sona/v1/create", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 400 {
		t.Errorf("Expected 400 status code got %v", w.Result())
	}
}

func TestIncidentUpdateWithValidToken(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})

	m := make(map[string]string, 1)
	m["Test"] = "Value"

	update := IncidentUpdate{"New State", "", "", m}
	_, token := user1.Authenticate("1234")
	body, _ := json.Marshal(update)

	r, _ := http.NewRequest("PUT", "/sona/v1/0/update", bytes.NewBuffer(body))
	r.Header.Set("X-Sona-Token", token.Token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 {
		t.Errorf("Expected 200 status code got %v", w.Result())
	}
}

func TestIncidentUpdateWithInvalidValidToken(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})

	m := make(map[string]string, 1)
	m["Test"] = "Value"

	update := IncidentUpdate{"New State", "", "", m}
	body, _ := json.Marshal(update)

	r, _ := http.NewRequest("PUT", "/sona/v1/0/update", bytes.NewBuffer(body))
	r.Header.Set("X-Sona-Token", "badValue")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 403 {
		t.Errorf("Expected 403 status code got %v", w.Result())
	}
}

func TestIncidentUpdateWithInvalidId(t *testing.T) {
	setup()
	body, _ := json.Marshal(IncidentUpdate{})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("PUT", "/sona/v1/badvalue/update", bytes.NewBuffer(body))
	r.Header.Set("X-Sona-Token", token.Token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 400 {
		t.Errorf("Expected 400 status code got %v", w.Result())
	}
}

func TestIncidentUpdateWithNonExistantId(t *testing.T) {
	setup()
	body, _ := json.Marshal(IncidentUpdate{})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("PUT", "/sona/v1/3/update", bytes.NewBuffer(body))
	r.Header.Set("X-Sona-Token", token.Token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 404 {
		t.Errorf("Expected 404 status code got %v", w.Result())
	}
}

func TestGetIncidentHandler(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("GET", "/sona/v1/0", nil)
	r.Header.Set("X-Sona-Token", token.Token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 {
		t.Errorf("Expected 200 status code got %v", w.Result())
	}

	var retVal Incident
	err := json.Unmarshal(w.Body.Bytes(), &retVal)
	if err != nil {
		t.Errorf("Failed to convert response %v error %v", w.Body, err)
	}

	if retVal.Id != 0 {
		t.Errorf("Expected incident 0 got %v", retVal.Id)
	}

	if retVal.Reporter != "Tester" {
		t.Errorf("Expected reporter Tester got %v", retVal.Reporter)
	}

	if retVal.Description != "Test" {
		t.Errorf("Expected description Test got %v", retVal.Description)
	}

	if retVal.State != "open" {
		t.Errorf("Expected state open got %v", retVal.State)
	}
}

func TestGetIncidentHandlerWithInvalidToken(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})

	r, _ := http.NewRequest("GET", "/sona/v1/0", nil)
	r.Header.Set("X-Sona-Token", "badToken")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 403 {
		t.Errorf("Expected 403 status code got %v", w.Result())
	}
}

func TestGetIncidentHandlerWithInvalidId(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("GET", "/sona/v1/zero", nil)
	r.Header.Set("X-Sona-Token", token.Token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 400 {
		t.Errorf("Expected 400 status code got %v", w.Result())
	}
}

func TestGetIncidentHandlerWithNonExistantId(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("GET", "/sona/v1/1", nil)
	r.Header.Set("X-Sona-Token", token.Token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 404 {
		t.Errorf("Expected 404 status code got %v", w.Result())
	}
}

func TestGetIncidentsHandler(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	incidentManager.AddIncident(&Incident{"Incident", 1, "Something", "Someone", "Closed", make(map[string]string, 0)})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("GET", "/sona/v1/incidents", nil)
	r.Header.Set("X-Sona-Token", token.Token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 {
		t.Errorf("Expected 200 status code got %v", w.Result())
	}

	var retVal []Incident
	err := json.Unmarshal(w.Body.Bytes(), &retVal)
	if err != nil {
		t.Errorf("Failed to convert response %v error %v", w.Body, err)
	}

	if len(retVal) != 2 {
		t.Errorf("Expected 2 incidents got %v", len(retVal))
	}

	if retVal[0].Id != 0 {
		t.Errorf("Expected incident 0 got %v", retVal[0].Id)
	}

	if retVal[0].Reporter != "Tester" {
		t.Errorf("Expected reporter Tester got %v", retVal[0].Reporter)
	}

	if retVal[0].Description != "Test" {
		t.Errorf("Expected description Test got %v", retVal[0].Description)
	}

	if retVal[0].State != "open" {
		t.Errorf("Expected state open got %v", retVal[0].State)
	}

	if retVal[1].Id != 1 {
		t.Errorf("Expected incident 0 got %v", retVal[1].Id)
	}

	if retVal[1].Reporter != "Someone" {
		t.Errorf("Expected reporter Someone got %v", retVal[1].Reporter)
	}

	if retVal[1].Description != "Something" {
		t.Errorf("Expected description Test got %v", retVal[1].Description)
	}

	if retVal[1].State != "Closed" {
		t.Errorf("Expected state Closed got %v", retVal[1].State)
	}
}

func TestGetIncidentsHandlerWithInvalidToken(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	incidentManager.AddIncident(&Incident{"Incident", 1, "Something", "Someone", "Closed", make(map[string]string, 0)})

	r, _ := http.NewRequest("GET", "/sona/v1/incidents", nil)
	r.Header.Set("X-Sona-Token", "badToken")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 403 {
		t.Errorf("Expected 403 status code got %v", w.Result())
	}
}

func TestGetAttachmentWithInvalidId(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("GET", "/sona/v1/zero/attachments", nil)
	r.Header.Set("X-Sona-Token", token.Token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 400 {
		t.Errorf("Expected 400 status code got %v", w.Result())
	}
}

func TestGetAttachmentsWithNoAttached(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("GET", "/sona/v1/0/attachments", nil)
	r.Header.Set("X-Sona-Token", token.Token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 {
		t.Errorf("Expected 200 status code got %v", w.Result())
	}

	var retVal []Attachment
	err := json.Unmarshal(w.Body.Bytes(), &retVal)
	if err != nil {
		t.Errorf("Failed to convert response %v error %v", w.Body, err)
	}

	if len(retVal) != 0 {
		t.Errorf("Expected 0 attachments but got %v", len(retVal))
	}
}

func TestGetAttachmentsWithAttached(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	incidentManager.AddAttachment(0, Attachment{"testfile.png", "2009-11-10T23:00:00Z"})
	incidentManager.AddAttachment(0, Attachment{"testfile2.jpg", "2009-10-10T23:00:00Z"})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("GET", "/sona/v1/0/attachments", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-Sona-Token", token.Token)

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 {
		t.Errorf("Expected 200 status code got %v", w.Result())
	}

	var retVal []Attachment
	err := json.Unmarshal(w.Body.Bytes(), &retVal)
	if err != nil {
		t.Errorf("Failed to convert response %v error %v", w.Body, err)
	}

	if len(retVal) != 2 {
		t.Errorf("Expected 2 attachments but got %v", len(retVal))
	}

	if retVal[0].FileName != "testfile.png" {
		t.Errorf("Expected first file to be testfile.png but got %v", retVal[0].FileName)
	}

	if retVal[1].FileName != "testfile2.jpg" {
		t.Errorf("Expected second file to be testfile2.jpg but got %v", retVal[0].FileName)
	}
}

func TestGetAttachmentsWithAttachedAndInvalidToken(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	incidentManager.AddAttachment(0, Attachment{"testfile.png", "2009-11-10T23:00:00Z"})
	incidentManager.AddAttachment(0, Attachment{"testfile2.jpg", "2009-10-10T23:00:00Z"})

	r, _ := http.NewRequest("GET", "/sona/v1/0/attachments", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-Sona-Token", "badToken")

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 403 {
		t.Errorf("Expected 403 status code got %v", w.Result())
	}
}

func TestUploadAttachmentWithInvalidId(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("POST", "/sona/v1/zero/attachment", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-Sona-Token", token.Token)

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 400 {
		t.Errorf("Expected 400 status code got %v", w.Result())
	}
}

func TestUploadAttachmentWithInvalidToken(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})

	r, _ := http.NewRequest("POST", "/sona/v1/0/attachment", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-Sona-Token", "badToken")

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 403 {
		t.Errorf("Expected 403 status code got %v", w.Result())
	}
}

func TestUploadAttachmentWithNonExistantId(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("POST", "/sona/v1/3/attachment", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-Sona-Token", token.Token)

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 404 {
		t.Errorf("Expected 404 status code got %v", w.Result())
	}
}

func TestDeleteAttachmentWithInvalidId(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("DELETE", "/sona/v1/zero/attachment/test.jpg", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-Sona-Token", token.Token)

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 400 {
		t.Errorf("Expected 400 status code got %v", w.Result())
	}
}

func TestDeleteAttachmentWithNonExistantIncidentId(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("DELETE", "/sona/v1/3/attachment/test.jpg", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-Sona-Token", token.Token)

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 404 {
		t.Errorf("Expected 404 status code got %v", w.Result())
	}
}

func TestDeleteAttachmentWithNonExistantAttachmentId(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	incidentManager.AddAttachment(0, Attachment{"somefile.png", "2009-11-10T23:00:00Z"})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("DELETE", "/sona/v1/0/attachment/test.jpg", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-Sona-Token", token.Token)

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 404 {
		t.Errorf("Expected 404 status code got %v", w.Result())
	}
}

func TestDeleteAttachment(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	incidentManager.AddAttachment(0, Attachment{"test.jpg", "2009-11-10T23:00:00Z"})
	_, token := user1.Authenticate("1234")

	r, _ := http.NewRequest("DELETE", "/sona/v1/0/attachment/test.jpg", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-Sona-Token", token.Token)

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 200 {
		t.Errorf("Expected 200 status code got %v", w.Result())
	}
}

func TestDeleteAttachmentWithInvalidToken(t *testing.T) {
	setup()
	incidentManager.AddIncident(&Incident{"Incident", 0, "Test", "Tester", "open", make(map[string]string, 0)})
	incidentManager.AddAttachment(0, Attachment{"test.jpg", "2009-11-10T23:00:00Z"})

	r, _ := http.NewRequest("DELETE", "/sona/v1/0/attachment/test.jpg", nil)
	w := httptest.NewRecorder()
	r.Header.Set("X-Sona-Token", "badToken")

	router.ServeHTTP(w, r)

	if w.Result().StatusCode != 403 {
		t.Errorf("Expected 403 status code got %v", w.Result())
	}
}
