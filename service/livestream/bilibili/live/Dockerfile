ARG base_image=justhumanz/go-simp:latest
FROM ${base_image} as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/livebili /livebili
RUN apk --no-cache add tzdata
CMD ["./livebili"]