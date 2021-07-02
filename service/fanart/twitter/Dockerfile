FROM justhumanz/go-simp as go-simp 
RUN apk --no-cache add tzdata

FROM alpine
COPY --from=go-simp /app/Go-Simp/bin/twitter_fanart /twitter_fanart
CMD ["./twitter_fanart"]