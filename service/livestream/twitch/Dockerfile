FROM go-simp
FROM alpine as twitch
COPY --from=go-simp /app/Go-Simp/bin/livetwitch /livetwitch
CMD ["./livetwitch"]