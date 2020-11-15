FROM golang:alpine

RUN apk add --update --no-cache git gcc build-base
RUN export CGO_ENABLED=1
RUN mkdir /app
WORKDIR /app
RUN git clone --single-branch --branch main https://github.com/JustHumanz/Go-Simp
WORKDIR /app/Go-Simp
RUN git checkout main
COPY config.toml /app/Go-Simp

#CMD ["go","run","."]
