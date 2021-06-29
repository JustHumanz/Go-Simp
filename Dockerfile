# build stage
FROM golang:alpine as build-stage

RUN apk add git
RUN export CGO_ENABLED=1
RUN mkdir /app
COPY . /app/Go-Simp
WORKDIR /app/Go-Simp
RUN rm -rf Img
WORKDIR /app/Go-Simp/bin

#fanart
RUN go build -o bilibili_fanart ../service/fanart/bilibili/
RUN go build -o pixiv_fanart ../service/fanart/pixiv/
RUN go build -o twitter_fanart ../service/fanart/twitter/


#frontend
RUN go build -o fe ../service/frontend/

#guild
RUN go build -o guild ../service/guild/

#live
RUN go build -o livebili ../service/livestream/bilibili/live
RUN go build -o spacebili ../service/livestream/bilibili/space
RUN go build -o liveyoutube ../service/livestream/youtube
RUN go build -o liveyoutube_counter ../service/livestream/youtube_counter
RUN go build -o livetwitch ../service/livestream/twitch

#migrate
RUN go build -o migrate ../service/migrate/

#pilot
RUN go build -o pilot ../service/pilot/

#api
RUN go build -o api ../service/rest-api/

#subscriber
RUN go build -o subscriber ../service/subscriber/

#utility
RUN go build -o utility ../service/utility/