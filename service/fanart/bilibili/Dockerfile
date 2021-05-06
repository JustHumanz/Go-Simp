FROM go-simp 
FROM alpine as bilibili_fanart
COPY --from=go-simp /app/Go-Simp/bin/bilibili_fanart /bilibili_fanart
CMD ["./bilibili_fanart"]
