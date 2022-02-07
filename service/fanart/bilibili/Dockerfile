FROM alpine
COPY --from=justhumanz/go-simp:${RELEASE_VERSION:-latest} /app/Go-Simp/bin/bilibili_fanart /bilibili_fanart
CMD ["./bilibili_fanart"]
