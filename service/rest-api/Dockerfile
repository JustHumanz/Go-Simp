ARG base_image=justhumanz/go-simp:latest
FROM ${base_image} as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/api /api
RUN apk --no-cache add tzdata
EXPOSE 2525
CMD ["./api"]