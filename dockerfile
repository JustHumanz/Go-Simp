FROM golang:alpine

RUN apk update && apk add git

RUN mkdir /apps/
COPY . /apps/
WORKDIR /apps

RUN go mod download
RUN go get -u github.com/JustHumanz/Go-simp

CMD ["go","run","."]