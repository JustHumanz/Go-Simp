FROM alpine
COPY --from=justhumanz/go-simp /app/Go-Simp/bin/api /api
RUN apk --no-cache add tzdata
EXPOSE 2525
CMD ["./api"]