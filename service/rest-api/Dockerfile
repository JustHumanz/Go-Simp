FROM go-simp
FROM alpine as api
COPY --from=go-simp /app/Go-Simp/bin/api /api
EXPOSE 2525
CMD ["./api"]