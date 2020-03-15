# Installation
Installation only works using docker images. However you can clone and build the server yourself if you would like to run it as not a docker image.

## Prerequisit note
In order to run Sona Server you will have to generate a configuration json file. Below is a very simple example file to get you up and running.

config.json
```json
{
    "adminConfig": {
        "emailAddress": "something@somewhere.com",
        "password": "itsasecret"
    },
    "managertype": 0,
    "filemanagertype": 0,
    "logging": {
        "enabled": false
    }
}
```

For more information on generating a configuration file see the following resources.

* [Configure File Manager](./ConfigureFileManager.md)
* [Configure Incident Manager](./ConfigureIncidentManager.md)
* [Configure Logging](./ConfigureLogging.md)
* [Configure Web Hooks](./ConfigureWebHooks.md)

## Docker
The easiest way to get up and running is to simply download the docker image and start it with your configuration file.

### Install Docker Image
`docker pull jeffriggle/sona-server:master`

### Start Docker Image
```shell
config=`cat config.json`
docker run -i -e CONFIG="$config" -p 8080:8080 jeffriggle/sona-server:master
```

## Build and Run

### Bash (Linux and Mac)
Building the server. In order to do this your machine will need git and golang installed on it. Once this is done you can clone this repo and build src. For an example of creating the binary yourself see [install.sh](../deploy/install.sh)

To run the server simply run the executable with a configuration file.

`/src ./Config.json`

### Powershell (Windows)
Building the server. In order to do this your machine will need git and golang installed on it. Once this is done you can clone this repo and build src. For an example of creating the binary yourself see [install.ps1](../deploy/install.ps1)

To run the server simply run the executable with a configuration file.

`src.exe ./Config.json`