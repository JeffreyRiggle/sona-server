# Configuration of the file manager
Currently sona server has a couple different options for file management (incident attachments). As of right now files can be stored on the local machine sona server is running on or in an [S3 bucket](https://aws.amazon.com/s3/).

## Selecting a file manager
The file manager is selected in the config file provided to sona. The valid options are

* 0 - local file system
* 1 - S3 bucket

```json
{
    "filemanagertype": 1
}
```

## Using the local file system.
By default sona uses the local file system to store files. When this option is used files will be stored in the home directory for the machine running the server.

You can override the storage path using the configuration file.

```json
{
    "filemanagertype": 0,
    "localfileconfig": "./some/path/here"
}
```

## Using S3
In order to use S3 a couple more configuration options need to be provided. In order to use the S3 manager the current assumption is that you will be running sona server from an [EC2 instance](https://aws.amazon.com/ec2/) with a role that allows for S3 access to the bucket you would like to store data in.

To use S3 you will also need to specific a region and name for your bucket.

```json
{
    "filemanagertype": 1,
    "s3config": {
        "region": "us-east-1",
        "bucket": "mybucket"
    }
}
```
