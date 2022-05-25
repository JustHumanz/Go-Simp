ARG BASE_IMAGE=justhumanz/go-simp:latest
FROM $BASE_IMAGE as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/subscriber_twitter /subscriber_twitter
RUN apk --no-cache add tzdata
CMD ["./subscriber_twitter"]