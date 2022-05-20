ARG BASE_IMAGE=justhumanz/go-simp:latest
FROM $BASE_IMAGE as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/liveyoutube_past_tracker /liveyoutube_past_tracker
RUN apk --no-cache add tzdata
CMD ["./liveyoutube_past_tracker"]