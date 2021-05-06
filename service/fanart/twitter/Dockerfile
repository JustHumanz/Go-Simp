FROM go-simp
FROM alpine as twitter_fanart
COPY --from=go-simp /app/Go-Simp/bin/twitter_fanart /twitter_fanart
CMD ["./twitter_fanart"]