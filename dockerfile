FROM golang:latest

RUN go get -u github.com/JustHumanz/Go-Simp
COPY config.toml $GOPATH/src/github.com/JustHumanz/Go-Simp/
WORKDIR $GOPATH/src/github.com/JustHumanz/Go-Simp
#CMD ["go","run","."]
