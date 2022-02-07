ARG base_image=justhumanz/go-simp:latest
FROM ${base_image} as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/migrate /migrate
RUN apk --no-cache add tzdata
CMD ["./migrate"]