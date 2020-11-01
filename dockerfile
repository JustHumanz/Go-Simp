FROM golang:alpine

RUN apk update && apk add git gcc build-base
RUN export CGO_ENABLED=1
RUN mkdir /app
WORKDIR /app
RUN go get ./...
RUN git clone https://github.com/JustHumanz/Go-Simp
COPY config.toml /app/Go-Simp

#CMD ["go","run","."]
