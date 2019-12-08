import requests
from assertpy import assert_that

# Initially basic test.
def test_create_user():
    return requests.post('http://localhost:8080/sona/v1/users', json={
        "emailAddress": "a@b.c",
        "userName": "TestUser",
        "firstName": "Test",
        "lastName": "User",
        "gender": "F",
        "password": "changeme"
        })

res = test_create_user()
assert_that(res.status_code).is_equal_to(201)