ARG BASE_IMAGE=justhumanz/go-simp:latest
FROM $BASE_IMAGE as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/subscriber_twitch /subscriber_twitch
RUN apk --no-cache add tzdata
CMD ["./subscriber_twitch"]