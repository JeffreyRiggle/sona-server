FROM ubuntu:latest

RUN apt-get update
RUN apt-get install jq -y
RUN apt-get install ca-certificates -y && rm -rf /var/cache/apk/*

RUN mkdir /app


RUN groupadd -r sona && useradd -r -s /bin/false -g sona sona
COPY ./src /app
COPY start.sh /app

RUN chown -R sona:sona /app
RUN chown -R sona:sona /home

RUN [ "chmod", "+x", "app/start.sh" ]

USER sona
ENTRYPOINT [ "app/start.sh" ]

EXPOSE 8080