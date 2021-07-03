FROM alpine
COPY --from=justhumanz/go-simp /app/Go-Simp/bin/liveyoutube /liveyoutube
RUN apk --no-cache add tzdata
CMD ["./liveyoutube"]