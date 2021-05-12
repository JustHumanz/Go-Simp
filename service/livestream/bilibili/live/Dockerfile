FROM go-simp
FROM alpine as livebili
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/livebili /livebili
CMD ["./livebili"]