FROM golang:alpine

RUN apk add git
RUN export CGO_ENABLED=1
RUN mkdir /app
COPY . /app/Go-Simp
WORKDIR /app/Go-Simp
RUN rm -rf Img
RUN go mod download
