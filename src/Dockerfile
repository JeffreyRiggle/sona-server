FROM ubuntu:latest

RUN apt-get update

COPY ./src /
COPY start.sh /

ENTRYPOINT [ "/start.sh" ]

EXPOSE 8080