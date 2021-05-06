FROM go-simp
FROM alpine as frontend
COPY --from=go-simp /app/Go-Simp/bin/frontend /frontend
CMD ["./frontend"]