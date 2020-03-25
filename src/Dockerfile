FROM ubuntu:latest

RUN apt-get update

RUN mkdir /app

RUN groupadd -r sona && useradd -r -s /bin/false -g sona sonauser
COPY ./src /app
COPY start.sh /app

RUN chown -R sonauser:sona /app/*

USER sonauser
ENTRYPOINT [ "/app/start.sh" ]

EXPOSE 8080