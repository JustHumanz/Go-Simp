ARG BASE_IMAGE=justhumanz/go-simp:latest
FROM $BASE_IMAGE as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/liveyoutube_upcoming_checker /liveyoutube_upcoming_checker
RUN apk --no-cache add tzdata
CMD ["./liveyoutube_upcoming_checker"]