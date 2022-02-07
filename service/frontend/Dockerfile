ARG base_image=justhumanz/go-simp:latest
FROM ${base_image} as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/fe /fe
RUN apk --no-cache add tzdata
CMD ["./fe"]