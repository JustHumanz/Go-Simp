FROM alpine
COPY --from=justhumanz/go-simp:${RELEASE_VERSION:-latest} /app/Go-Simp/bin/pixiv_fanart /pixiv_fanart
CMD ["./pixiv_fanart"]
