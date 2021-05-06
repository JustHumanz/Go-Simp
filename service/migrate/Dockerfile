FROM go-simp
FROM alpine as migrate
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/migrate /migrate
CMD ["./migrate"]