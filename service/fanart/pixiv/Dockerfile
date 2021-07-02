FROM justhumanz/go-simp as go-simp
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/pixiv_fanart /pixiv_fanart
CMD ["./pixiv_fanart"]
