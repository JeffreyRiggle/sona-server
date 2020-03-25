FROM ubuntu:latest

RUN apt-get update

RUN mkdir /app

RUN groupadd -r sona && useradd -r -s /bin/false -g sona sona
COPY ./src /app
COPY start.sh /app

RUN chown -R sona:sona /app
RUN [ "chmod", "+x", "app/start.sh" ]

USER sona
ENTRYPOINT [ "app/start.sh" ]

EXPOSE 8080