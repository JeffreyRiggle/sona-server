import requests
import json
import testrunner
import time
from assertpy import assert_that

incAdminToken = ''
restrictedToken = ''
webhookIncident = {}

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
        "emailAddress": ures.json().get("emailAddress"),
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
        "emailAddress": 'something@somewhere.com',
        "password": "itsasecret"
        })
    incAdminToken = res.json().get("token")

    res = requests.put('http://localhost:8080/sona/v1/incidents/0', headers={'X-Sona-Token': incAdminToken}, json={
        "state": "In Progress",
        })
    
    assert_that(res.status_code).is_equal_to(200)

def test_attach_without_auth():
    files = {'uploadfile': open('attach.txt', 'rb')} 
    res = requests.post('http://localhost:8080/sona/v1/incidents/0/attachment', files=files)
    assert_that(res.status_code).is_equal_to(403)

def test_attach_without_permission():
    global restrictedToken

    files = {'uploadfile': open('attach.txt', 'rb')} 
    res = requests.post('http://localhost:8080/sona/v1/incidents/0/attachment', headers={'X-Sona-Token': restrictedToken}, files=files)
    assert_that(res.status_code).is_equal_to(401)

def test_attach():
    global incAdminToken

    files = {'uploadfile': open('attach.txt', 'rb')} 
    res = requests.post('http://localhost:8080/sona/v1/incidents/0/attachment', headers={'X-Sona-Token': incAdminToken}, files=files)
    
    assert_that(res.status_code).is_equal_to(200)

def test_get_attachments_without_auth():
        res = requests.get('http://localhost:8080/sona/v1/incidents/0/attachments')
        assert_that(res.status_code).is_equal_to(403)

def test_get_attachments_without_permission():
        global restrictedToken
        
        res = requests.get('http://localhost:8080/sona/v1/incidents/0/attachments', headers={'X-Sona-Token': restrictedToken})
        assert_that(res.status_code).is_equal_to(401)

def test_get_attachments():
    global incAdminToken
        
    res = requests.get('http://localhost:8080/sona/v1/incidents/0/attachments', headers={'X-Sona-Token': incAdminToken})
    assert_that(res.status_code).is_equal_to(200)

    content = res.json()
    assert_that(len(content)).is_equal_to(1)
    
    attach = content[0]
    assert_that(attach.get("filename")).is_equal_to("attach.txt")

def test_download_attachment_without_auth():
    res = requests.get('http://localhost:8080/sona/v1/incidents/0/attachment/attach.txt')
    assert_that(res.status_code).is_equal_to(403)

def test_download_attachment_without_permission():
    global restrictedToken
    res = requests.get('http://localhost:8080/sona/v1/incidents/0/attachment/attach.txt', headers={'X-Sona-Token': restrictedToken})
    assert_that(res.status_code).is_equal_to(401)

def test_download_attachment():
    global incAdminToken
    res = requests.get('http://localhost:8080/sona/v1/incidents/0/attachment/attach.txt', headers={'X-Sona-Token': incAdminToken})
    assert_that(res.status_code).is_equal_to(200)

def test_delete_attachment_without_auth():
    res = requests.delete('http://localhost:8080/sona/v1/incidents/0/attachment/attach.txt')
    assert_that(res.status_code).is_equal_to(403)

def test_delete_attachment_without_permission():
    global restrictedToken

    res = requests.delete('http://localhost:8080/sona/v1/incidents/0/attachment/attach.txt', headers={'X-Sona-Token': restrictedToken})
    assert_that(res.status_code).is_equal_to(401)

def test_delete_attachment():
    global incAdminToken

    res = requests.delete('http://localhost:8080/sona/v1/incidents/0/attachment/attach.txt', headers={'X-Sona-Token': incAdminToken})
    assert_that(res.status_code).is_equal_to(200)

def test_get_incident_without_auth():
    res = requests.get('http://localhost:8080/sona/v1/incidents/0')
    assert_that(res.status_code).is_equal_to(403)

def test_get_incident_without_permission():
    global restrictedToken

    res = requests.get('http://localhost:8080/sona/v1/incidents/0/', headers={'X-Sona-Token': restrictedToken})
    assert_that(res.status_code).is_equal_to(401)

def test_get_incident():
    global incAdminToken

    res = requests.get('http://localhost:8080/sona/v1/incidents/0', headers={'X-Sona-Token': incAdminToken})
    assert_that(res.status_code).is_equal_to(200)

    inc = res.json()

    assert_that(inc.get("id")).is_equal_to(0)
    assert_that(inc.get("description")).is_equal_to("Something is wrong")
    assert_that(inc.get("reporter")).is_equal_to("TestUser")
    assert_that(inc.get("state")).is_equal_to("In Progress")

def test_get_filtered_incidents():
    global incAdminToken

    res = requests.post('http://localhost:8080/sona/v1/incidents', json={
        "description": "Something else is wrong",
        "reporter": "Steve",
        "state": "open",
        "attributes": {
            "Foo": "Bar"
        }
        })
    
    assert_that(res.status_code).is_equal_to(201)
    res = requests.post('http://localhost:8080/sona/v1/incidents', json={
        "description": "Something new is wrong",
        "reporter": "Jill",
        "state": "open",
        "attributes": {
            "Foo": "Bar"
        }
        })
    assert_that(res.status_code).is_equal_to(201)

    incFilter = '?filter=%7B%22complexfilters%22%3A%5B%7B%22filters%22%3A%5B%7B%22property%22%3A%22Reporter%22%2C%22comparison%22%3A%22equals%22%2C%22value%22%3A%22Jill%22%7D%5D%2C%22junction%22%3A%22and%22%7D%5D%2C%22union%22%3A%22and%22%7D'
    req = 'http://localhost:8080/sona/v1/incidents' + incFilter
    res = requests.get(req, headers={'X-Sona-Token': incAdminToken})
    
    assert_that(res.status_code).is_equal_to(200)
    content = res.json()

    assert_that(len(content)).is_equal_to(1)
    assert_that(content[0].get("description")).is_equal_to("Something new is wrong")

def test_add_incident_hook():
    global webhookIncident

    res = requests.delete('http://localhost:5000/calls')
    assert_that(res.status_code).is_equal_to(200)

    res = requests.post('http://localhost:8080/sona/v1/incidents', json={
        "description": "Hook Testing",
        "reporter": "captin",
        "state": "open",
        "attributes": {
            "Foo": "Bar"
        }
    })
    assert_that(res.status_code).is_equal_to(201)

    webhookIncident = res.json()

    time.sleep(5)
    hooks = requests.get('http://localhost:5000/calls')
    addedCalls = hooks.json().get("incidentAdded")
    assert_that(len(addedCalls)).is_equal_to(1)
    assert_that(addedCalls[0].get("body")).is_equal_to("New Incident Created by captin with description Hook Testing.")
    assert_that(addedCalls[0].get("incident")).is_equal_to(str(webhookIncident.get("id")))
    assert_that(addedCalls[0].get("subject")).is_equal_to("Incident Created")
    assert_that(addedCalls[0].get("to")).is_equal_to("foobar@email.com")

def test_update_incident_hook():
    global incAdminToken
    global webhookIncident

    res = requests.put('http://localhost:8080/sona/v1/incidents/' + str(webhookIncident.get("id")), headers={'X-Sona-Token': incAdminToken}, json={
        "state": "In Progress",
        })
    
    assert_that(res.status_code).is_equal_to(200)

    time.sleep(5)
    hooks = requests.get('http://localhost:5000/calls')
    updatedCalls = hooks.json().get("incidentUpdated")
    assert_that(len(updatedCalls)).is_equal_to(1)
    assert_that(updatedCalls[0].get("body")).is_equal_to("Incident " + str(webhookIncident.get("id")) + " updated with description  and state In Progress.")
    assert_that(updatedCalls[0].get("incident")).is_equal_to(str(webhookIncident.get("id")))
    assert_that(updatedCalls[0].get("subject")).is_equal_to("Incident Updated")
    assert_that(updatedCalls[0].get("to")).is_equal_to("foobar@email.com")

def test_attach_incident_hook():
    global incAdminToken
    global webhookIncident

    files = {'uploadfile': open('attach.txt', 'rb')} 
    res = requests.post('http://localhost:8080/sona/v1/incidents/' + str(webhookIncident.get("id")) + '/attachment', headers={'X-Sona-Token': incAdminToken}, files=files)
    
    assert_that(res.status_code).is_equal_to(200)

    time.sleep(5)
    hooks = requests.get('http://localhost:5000/calls')
    attachedCalls = hooks.json().get("incidentAttached")
    assert_that(len(attachedCalls)).is_equal_to(1)
    assert_that(attachedCalls[0].get("body")).is_equal_to("attach.txt Has been attached to Incident " + str(webhookIncident.get("id")))
    assert_that(attachedCalls[0].get("incident")).is_equal_to(str(webhookIncident.get("id")))
    assert_that(attachedCalls[0].get("subject")).is_equal_to("Attachment added")
    assert_that(attachedCalls[0].get("to")).is_equal_to("foobar@email.com")

def setup():
    testrunner.addTest("Create Incident", test_create_incident)
    testrunner.addTest("Update Incident without auth", test_update_incident_without_auth)
    testrunner.addTest("Update Incident without permission", test_update_incident_without_permission)
    testrunner.addTest("Update Incident", test_update_incident)
    testrunner.addTest("Attach to Incident without auth", test_attach_without_auth)
    testrunner.addTest("Attach to Incident without permission", test_attach_without_permission)
    testrunner.addTest("Attach to Incident", test_attach)
    testrunner.addTest("Get attachments without auth", test_get_attachments_without_auth)
    testrunner.addTest("Get attachments without permission", test_get_attachments_without_permission)
    testrunner.addTest("Get attachments", test_get_attachments)
    testrunner.addTest("Test download attachment without auth", test_download_attachment_without_auth)
    testrunner.addTest("Test download attachment without permission", test_download_attachment_without_permission)
    testrunner.addTest("Test download attachment", test_download_attachment)
    testrunner.addTest("Test delete attachment without auth", test_delete_attachment_without_auth)
    testrunner.addTest("Test delete attachment without permission", test_delete_attachment_without_permission)
    testrunner.addTest("Test delete attachment", test_delete_attachment)
    testrunner.addTest("Test get incident without auth", test_get_incident_without_auth)
    testrunner.addTest("Test get incident without permission", test_get_incident_without_permission)
    testrunner.addTest("Test get incident", test_get_incident)
    testrunner.addTest("Test get filtered incident", test_get_filtered_incidents)
    testrunner.addTest("Test incident added hook", test_add_incident_hook)
    testrunner.addTest("Test incident updated hook", test_update_incident_hook)
    testrunner.addTest("Test incident attached hook", test_attach_incident_hook)