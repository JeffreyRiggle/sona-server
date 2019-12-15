import requests
import testrunner
from assertpy import assert_that

incAdminToken = ''
restrictedToken = ''

def test_create_incident():
    res = requests.post('http://localhost:8080/sona/v1/incidents', json={
        "description": "Something is wrong",
        "reporter": "TestUser",
        "state": "open",
        "attributes": {
            "Test": "Value"
        }
        })
    assert_that(res.status_code).is_equal_to(201)

    inc1 = res.json()
    assert_that(inc1.get("id")).is_equal_to(0)
    assert_that(inc1.get("description")).is_equal_to("Something is wrong")
    assert_that(inc1.get("reporter")).is_equal_to("TestUser")
    assert_that(inc1.get("state")).is_equal_to("open")

    attrs = inc1.get("attributes")
    assert_that(attrs.get("Test")).is_equal_to("Value")

def test_update_incident_without_auth():
    res = requests.put('http://localhost:8080/sona/v1/incidents/0', json={
        "state": "In Progress",
        })
    assert_that(res.status_code).is_equal_to(403)

def test_update_incident_without_permission():
    global restrictedToken
    ures = requests.post('http://localhost:8080/sona/v1/users', json={
        "emailAddress": "email@address.com",
        "userName": "IncTest",
        "firstName": "Incident",
        "lastName": "User",
        "gender": "F",
        "password": "1234"
        })

    ares = requests.post('http://localhost:8080/sona/v1/authenticate', json={
        "id": ures.json().get("id"),
        "password": "1234"
        })

    restrictedToken = ares.json().get("token")

    res = requests.put('http://localhost:8080/sona/v1/incidents/0', headers={'X-Sona-Token': restrictedToken}, json={
        "state": "In Progress",
        })
    assert_that(res.status_code).is_equal_to(401)

def test_update_incident():
    global incAdminToken

    res = requests.post('http://localhost:8080/sona/v1/authenticate', json={
        "id": 0,
        "password": "admin"
        })
    incAdminToken = res.json().get("token")

    res = requests.put('http://localhost:8080/sona/v1/incidents/0', headers={'X-Sona-Token': incAdminToken}, json={
        "state": "In Progress",
        })
    
    assert_that(res.status_code).is_equal_to(200)

def test_attach_without_auth():
    res = requests.post('http://localhost:8080/sona/v1/incidents/0/attachment')
    assert_that(res.status_code).is_equal_to(403)

def test_attach_without_permission():
    global restrictedToken

    res = requests.post('http://localhost:8080/sona/v1/incidents/0/attachment', headers={'X-Sona-Token': restrictedToken})
    assert_that(res.status_code).is_equal_to(401)

def test_attach():
    global incAdminToken

    files = {'uploadfile': open('attach.txt', 'rb')} 
    res = requests.post('http://localhost:8080/sona/v1/incidents/0/attachment', headers={'X-Sona-Token': incAdminToken}, files=files)
    
    assert_that(res.status_code).is_equal_to(200)

def setup():
    testrunner.addTest("Create Incident", test_create_incident)
    testrunner.addTest("Update Incident without auth", test_update_incident_without_auth)
    testrunner.addTest("Update Incident without permission", test_update_incident_without_permission)
    testrunner.addTest("Update Incident", test_update_incident)
    testrunner.addTest("Attach to Incident without auth", test_attach_without_auth)
    testrunner.addTest("Attach to Incident without permission", test_attach_without_permission)
    testrunner.addTest("Attach to Incident", test_attach)