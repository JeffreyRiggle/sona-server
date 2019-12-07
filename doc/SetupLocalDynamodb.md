# Setting up Sona to use local dynamodb instance

## Install the dynamodb docker image
First you will need to install aws's dynamodb image for docker `sudo docker pull amazon/dynamodb-local`

## Start the dynamodb docker image
Next you will need to actually run the docker image `sudo docker run -p 8000:8000 amazon/dynamodb-local`

## Configure the dynamodb docker image with an access key
Navigate to http://localhost:8000/shell. Once on this page press the gear icon and set the access key. You can either use your own instance of AWS to generate a key using IAM or you can use the following Access Key `AKIAIOSFODNN7EXAMPLE` (note this is not very safe to do).

## Setting up Sona
In order to run against a local dynamodb instance you will need to do a couple things to the sona server. First you will need to create an aws secret file on the machine sona is running on. To do this create a file at the location `~/.aws/credentials`. In this file you will need to store your credentails this should look something like this.

```
[default]
aws_access_key_id = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
```

Finally you will have to update your configuration file to be able to run as dynamodb

```json
{
    "managertype": 1,
    "dynamodb": {
        "region": "us-east-1",
        "endpoint": "http://localhost:8000"
    }
    //...
}
```