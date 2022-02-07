FROM alpine
COPY --from=justhumanz/go-simp:${RELEASE_VERSION:-latest} /app/Go-Simp/bin/migrate /migrate
RUN apk --no-cache add tzdata
CMD ["./migrate"]