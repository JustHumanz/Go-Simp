ARG base_image=justhumanz/go-simp:latest
FROM ${base_image} as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/bilibili_fanart /bilibili_fanart
CMD ["./bilibili_fanart"]
