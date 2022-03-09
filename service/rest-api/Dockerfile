ARG BASE_IMAGE=justhumanz/go-simp:latest
FROM $BASE_IMAGE as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/api /api
RUN apk --no-cache add tzdata
EXPOSE 2525
CMD ["./api"]