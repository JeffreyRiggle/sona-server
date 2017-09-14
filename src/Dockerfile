FROM golang

ADD . /go/src/bitbucket.com/jeffreyriggle/sona

WORKDIR /go/src/bitbucket.com/jeffreyriggle/sona

RUN go get && go install bitbucket.com/jeffreyriggle/sona

ENTRYPOINT /go/bin/sona

EXPOSE 8080