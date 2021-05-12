FROM go-simp 
FROM alpine as bilibili_fanart
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/bilibili_fanart /bilibili_fanart
CMD ["./bilibili_fanart"]
