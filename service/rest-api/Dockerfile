FROM go-simp
FROM alpine as api
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/api /api
EXPOSE 2525
CMD ["./api"]