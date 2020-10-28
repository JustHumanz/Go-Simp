FROM golang:alpine

RUN apk update && apk add git gcc build-base
RUN export CGO_ENABLED=1
RUN go get -u github.com/JustHumanz/Go-Simp
COPY config.toml $GOPATH/src/github.com/JustHumanz/Go-Simp/
WORKDIR $GOPATH/src/github.com/JustHumanz/Go-Simp
RUN git pull
#CMD ["go","run","."]
