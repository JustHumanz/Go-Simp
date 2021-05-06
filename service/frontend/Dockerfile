FROM go-simp
FROM alpine as frontend
COPY --from=go-simp /app/Go-Simp/bin/fe /fe
CMD ["./fe"]