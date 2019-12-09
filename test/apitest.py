import requests
from assertpy import assert_that

auth = ''

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
    global auth
    res = requests.post('http://localhost:8080/sona/v1/authenticate', json={
        "id": 1,
        "password": "changeme"
        })
    assert_that(res.status_code).is_equal_to(200)
    auth = res.json().get("token")
    assert_that(len(auth)).is_greater_than(1)

def test_update_user():
    global auth
    res = requests.put('http://localhost:8080/sona/v1/users/1', headers={"X-Sona-Token": auth}, json={
        "gender": "M",
        })
    assert_that(res.status_code).is_equal_to(200)

test_create_user()
test_change_user_without_auth()
test_auth_user()
test_update_user()
# TODO update permissions
# TODO get other user
# TODO change password
# TODO test delete user