ARG base_image=justhumanz/go-simp:latest
FROM ${base_image} as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/pixiv_fanart /pixiv_fanart
CMD ["./pixiv_fanart"]
