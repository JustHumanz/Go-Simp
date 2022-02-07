FROM alpine
COPY --from=justhumanz/go-simp:${RELEASE_VERSION:-latest} /app/Go-Simp/bin/utility /utility
CMD ["./utility"]