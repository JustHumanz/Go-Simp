ARG base_image=justhumanz/go-simp:latest
FROM ${base_image} as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/pilot /pilot
EXPOSE 9000
CMD ["./pilot"]