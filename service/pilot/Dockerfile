ARG BASE_IMAGE=justhumanz/go-simp:latest
FROM $BASE_IMAGE as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/pilot /pilot
EXPOSE 9000 8080 8181
CMD ["./pilot"]