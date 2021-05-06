FROM go-simp
FROM alpine as frontend
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/fe /fe
CMD ["./fe"]