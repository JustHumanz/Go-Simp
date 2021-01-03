FROM golang:alpine

RUN apk add git
RUN export CGO_ENABLED=1
RUN mkdir /app
COPY . /app/Go-Simp
WORKDIR /app/Go-Simp
COPY config.toml /app/Go-Simp

#CMD ["go","run","."]
