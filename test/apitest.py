import requests
from assertpy import assert_that

userToken = ''
adminToken = ''

def test_create_user():
    res = requests.post('http://localhost:8080/sona/v1/users', json={
        "emailAddress": "a@b.c",
        "userName": "TestUser",
        "firstName": "Test",
        "lastName": "User",
        "gender": "F",
        "password": "changeme"
        })
    assert_that(res.status_code).is_equal_to(201)

    user1 = res.json()
    assert_that(user1.get("emailAddress")).is_equal_to("a@b.c")
    assert_that(user1.get("userName")).is_equal_to("TestUser")
    assert_that(user1.get("firstName")).is_equal_to("Test")
    assert_that(user1.get("lastName")).is_equal_to("User")
    assert_that(user1.get("gender")).is_equal_to("F")
    assert_that(user1.get("password")).is_equal_to(None)
    assert_that(len(user1.get("permissions"))).is_equal_to(0)
    assert_that(user1.get("id")).is_equal_to(1)

def test_change_user_without_auth():
    res = requests.put('http://localhost:8080/sona/v1/users/1', json={
        "gender": "M",
        })
    assert_that(res.status_code).is_equal_to(403)

def test_auth_user():
    global userToken
    res = requests.post('http://localhost:8080/sona/v1/authenticate', json={
        "id": 1,
        "password": "changeme"
        })
    assert_that(res.status_code).is_equal_to(200)
    userToken = res.json().get("token")
    assert_that(len(userToken)).is_greater_than(1)

def test_update_user():
    global userToken
    res = requests.put('http://localhost:8080/sona/v1/users/1', headers={"X-Sona-Token": userToken}, json={
        "gender": "M",
        })
    assert_that(res.status_code).is_equal_to(200)

def test_update_permissions_without_auth():
    res = requests.put('http://localhost:8080/sona/v1/users/1/permissions', json=["*"])
    assert_that(res.status_code).is_equal_to(403)

def test_update_permissions_with_invalid_auth():
    global userToken
    res = requests.put('http://localhost:8080/sona/v1/users/1/permissions', headers={"X-Sona-Token": userToken}, json=["*"])
    assert_that(res.status_code).is_equal_to(401)

def test_update_permissions():
    global adminToken
    res = requests.post('http://localhost:8080/sona/v1/authenticate', json={
        "id": 0,
        "password": "admin"
        })
    adminToken = res.json().get("token")

    res = requests.put('http://localhost:8080/sona/v1/users/1/permissions', headers={"X-Sona-Token": adminToken}, json=["view-incident"])
    assert_that(res.status_code).is_equal_to(200)

def test_get_user_without_auth():
    res = requests.get('http://localhost:8080/sona/v1/users/0')
    assert_that(res.status_code).is_equal_to(403)

def test_get_user_without_permission():
    res = requests.get('http://localhost:8080/sona/v1/users/0', headers={'X-Sona-Token': userToken})
    assert_that(res.status_code).is_equal_to(401)

def test_get_user():
    global adminToken

    res = requests.get('http://localhost:8080/sona/v1/users/1', headers={'X-Sona-Token': adminToken})
    assert_that(res.status_code).is_equal_to(200)
    content = res.json()
    assert_that(content.get("id")).is_equal_to(1)
    assert_that(content.get("emailAddress")).is_equal_to("a@b.c")
    assert_that(content.get("userName")).is_equal_to("TestUser")
    assert_that(content.get("firstName")).is_equal_to("Test")
    assert_that(content.get("lastName")).is_equal_to("User")
    assert_that(content.get("gender")).is_equal_to("M")
    assert_that(content.get("password")).is_equal_to(None)
    assert_that(len(content.get("permissions"))).is_equal_to(1)

def change_password_without_token():
    res = requests.put('http://localhost:8080/sona/v1/users/1/authentication', json={
        "oldPassword": "changeme",
        "newPassword": "1234"
    })
    assert_that(res.status_code).is_equal_to(403)

def change_password_with_incorrect_old_password():
    global userToken
    res = requests.put('http://localhost:8080/sona/v1/users/1/authentication', headers={'X-Sona-Token': userToken}, json={
        "oldPassword": "foobar",
        "newPassword": "1234"
    })
    assert_that(res.status_code).is_equal_to(403)

def change_password():
    global userToken
    res = requests.put('http://localhost:8080/sona/v1/users/1/authentication', headers={'X-Sona-Token': userToken}, json={
        "oldPassword": "changeme",
        "newPassword": "1234"
    })
    assert_that(res.status_code).is_equal_to(200)

def delete_user_without_token():
    res = requests.delete('http://localhost:8080/sona/v1/users/1')
    assert_that(res.status_code).is_equal_to(403)

def delete_user_without_permission():
    global userToken
    res = requests.delete('http://localhost:8080/sona/v1/users/0', headers={'X-Sona-Token': userToken})
    assert_that(res.status_code).is_equal_to(401)

def delete_user():
    global adminToken
    res = requests.delete('http://localhost:8080/sona/v1/users/1', headers={'X-Sona-Token': adminToken})
    assert_that(res.status_code).is_equal_to(200)

test_create_user()
test_change_user_without_auth()
test_auth_user()
test_update_user()
test_update_permissions_without_auth()
test_update_permissions_with_invalid_auth()
test_update_permissions()
test_get_user_without_auth()
test_get_user_without_permission()
test_get_user()
change_password_without_token()
change_password_with_incorrect_old_password()
change_password()
delete_user_without_token()
delete_user_without_permission()
delete_user()