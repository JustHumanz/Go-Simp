FROM golang:alpine

RUN apk add --update --no-cache git
RUN export CGO_ENABLED=1
RUN mkdir /app
WORKDIR /app
RUN git clone https://github.com/JustHumanz/Go-Simp
WORKDIR /app/Go-Simp
COPY config.toml /app/Go-Simp

#CMD ["go","run","."]
