# Setting up Sona to use mysql database

## Create the database
This tutorial will be using mysql on the local machine using a docker image. In order to do this you must first download [docker](docker). Once this is done to get the correct docker image run the following `sudo docker pull mysql/mysql-server`.

Now that we have a docker image available lets create a mysql server run the following command `sudo docker run --name=mysql1 -d -p 3306:3306 mysql/mysql-server`.

Wait for the image to be healthy. `sudo docker ps`.

Get the current password `sudo docker logs mysql1 2>&1 | grep GENERATED`.

Log into the server `sudo docker exec -it mysql2 mysql -uroot -p` enter the password.

Change the password to something else. `ALTER USER 'root'@'localhost' IDENTIFIED BY 'password';`

Add permissions for outside connections `CREATE USER 'root'@'%' IDENTIFIED BY '1234';` followed by `GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' WITH GRANT OPTION;`

Test your connection to the database `mysql -h localhost --protocol=tcp -uroot --password 1234 --port 3306`

Create a database `CREATE DATABASE sona;`

update the config

```json
"mysql": {
    "username": "root",
    "password": "1234",
    "host": "127.0.0.1",
    "port": "3306",
    "dbname": "sona"
}
```