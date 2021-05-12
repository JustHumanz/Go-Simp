FROM go-simp
FROM alpine as youtube
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/liveyoutube /liveyoutube
CMD ["./liveyoutube"]