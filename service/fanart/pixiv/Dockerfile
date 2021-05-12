FROM go-simp
FROM alpine as pixiv_fanart
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/pixiv_fanart /pixiv_fanart
CMD ["./pixiv_fanart"]
