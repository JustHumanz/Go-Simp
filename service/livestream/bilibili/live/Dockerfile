FROM alpine
COPY --from=justhumanz/go-simp:${RELEASE_VERSION:-latest} /app/Go-Simp/bin/livebili /livebili
RUN apk --no-cache add tzdata
CMD ["./livebili"]