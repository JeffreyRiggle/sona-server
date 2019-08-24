# Incident Management
Currently sona server has a couple different options for incident management. As of right now incidents can be stored on the local machine sona server is running on, aws's [dynamodb](https://aws.amazon.com/dynamodb/), [MySQL](https://www.mysql.com/) or in google's [Datastore](https://cloud.google.com/datastore/)

## Selecting an Incident Manager
The file incident manager is selected in the config file provided to sona. The valid options are

* 0 - runtime
* 1 - dynamodb
* 2 - MySQL
* 3 - Datastore

```json
{
    "incidentmanagertype": 0
}
```

## Using runtime incident manager

> Warning runtime is volitile and nothing will end up stored after sona server shuts down.

To use the runtime manager use the following json.
```json
{
    "incidentmanagertype": 0
}
```

All incidents will be purged on restart. This is probably only useful for getting a feel for how this service works with the need to setup a database.

## Using dynamodb
In order to use dynamodb a couple more configuration options need to be provided. In order to use the dynamodb the current assumption is that you will be running sona server from an [EC2 instance](https://aws.amazon.com/ec2/) with a role that allows for dynamodb access.

To use dynamodb you will need to use json similar to this.

```json
{
    "incidentmanagertype": 1,
    "dynamodb": {
        "region": "us-east-1"
    }
}
```

## Using MySQL
In order to use MySQL sona server will have to be able to authenticate with an MySQL server.

There is additional configuration required for this.
```json
{
    "incidentmanagertype": 2,
    "mysql": {
        "username": "admin",
        "password": "1234",
        "host": "127.0.0.1",
        "port": "3306",
        "dbname": "sona"
    }
}
```

## Using Datastore
In order to use Datastore you will need to create a datastore and you will need to generate a token.json file for access to that datastore. Once this is done you will need to provide this information to sona server.

```json
{
    "incidentmanagertype": 3,
    "datastore": {
        "projectname": "sona-181109",
        "authfile": "creds.json"
    }
}
```