ARG BASE_IMAGE=justhumanz/go-simp:latest
FROM $BASE_IMAGE as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/migrate /migrate
RUN apk --no-cache add tzdata
CMD ["./migrate"]