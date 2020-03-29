# Security
Sona has the ability to handle web traffic using http or https. By default Sona will handle traffic with http and custom configuration is required to use https. This will explain how to use https.

## Passing Configuration through to docker
In the case of using the docker image the https certificate and key will need to be passed in to the container using the configuration file.

In this case you will need to base64 encode your certificate and key and use those values in the configuration file.

```json
{
    "securityConfig": {
        "key": "base64 encoded key here",
        "cert": "base64 encoded cert here"
    }
}
```

## Using https in a non-docker image.
In the case of a non-docker image you can just specify the path to the key and certificate on the machine.

```json
{
    "securityConfig": {
        "key": "path to key",
        "cert": "path to certificate"
    }
}
```