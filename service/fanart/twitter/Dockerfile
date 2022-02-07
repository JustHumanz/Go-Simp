ARG base_image=justhumanz/go-simp:latest
FROM ${base_image} as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/twitter_fanart /twitter_fanart
CMD ["./twitter_fanart"]