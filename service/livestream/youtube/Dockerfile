FROM go-simp
FROM alpine as youtube
COPY --from=go-simp /app/Go-Simp/bin/liveyoutube /liveyoutube
CMD ["./liveyoutube"]