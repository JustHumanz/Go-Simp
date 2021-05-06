FROM go-simp
FROM alpine as guild
COPY --from=go-simp /app/Go-Simp/bin/guild /guild
CMD ["./guild"]