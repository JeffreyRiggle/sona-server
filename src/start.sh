#!/bin/bash

echo $CONFIG > app/start.json

jq .securityConfig.cert app/start.json | cut -d "\"" -f 2 > app/temp

if [ -s app/temp ]
then
    echo "Found cert"
    base64 -d app/temp > app/server.cert
    jq '.securityConfig.cert="app/server.cert"' app/start.json > app/temp
    echo `cat app/temp` > app/start.json
else
    echo "No certificate provided"
fi

jq .securityConfig.key app/start.json | cut -d "\"" -f 2 > app/temp

if [ -s app/temp ]
then
    echo "Found key"
    base64 -d app/temp > app/server.key
    jq '.securityConfig.key="app/server.key"' app/start.json > app/temp
    echo `cat app/temp` > app/start.json
else
    echo "No key provided"
fi

rm app/temp
app/src app/start.json