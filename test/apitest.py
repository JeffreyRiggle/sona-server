import requests
from assertpy import assert_that

# Initially basic test.
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

    content = res.json()
    assert_that(content.get("emailAddress")).is_equal_to("a@b.c")
    assert_that(content.get("userName")).is_equal_to("TestUser")
    assert_that(content.get("firstName")).is_equal_to("Test")
    assert_that(content.get("lastName")).is_equal_to("User")
    assert_that(content.get("gender")).is_equal_to("F")
    assert_that(content.get("password")).is_equal_to(None)
    assert_that(len(content.get("permissions"))).is_equal_to(0)

test_create_user()