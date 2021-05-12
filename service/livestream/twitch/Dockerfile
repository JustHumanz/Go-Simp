FROM go-simp
FROM alpine as twitch
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/livetwitch /livetwitch
CMD ["./livetwitch"]