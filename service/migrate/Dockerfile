FROM go-simp
FROM alpine as migrate
COPY --from=go-simp /app/Go-Simp/bin/migrate /migrate
CMD ["./migrate"]