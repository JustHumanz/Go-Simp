FROM go-simp
FROM alpine as youtube-counter
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/liveyoutube_counter /liveyoutube_counter
CMD ["./liveyoutube_counter"]