FROM golang:alpine

RUN apk update && apk add git
RUN go get -u github.com/JustHumanz/Go-Simp
COPY config/config.json $GOPATH/src/github.com/JustHumanz/Go-Simp/config

WORKDIR $GOPATH/src/github.com/JustHumanz/Go-Simp
CMD ["go","run","."]
