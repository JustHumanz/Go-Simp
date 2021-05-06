FROM go-simp
FROM alpine as livebili
COPY --from=go-simp /app/Go-Simp/bin/livebili /livebili
CMD ["./livebili"]