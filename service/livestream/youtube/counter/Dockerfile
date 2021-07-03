FROM alpine
COPY --from=justhumanz/go-simp /app/Go-Simp/bin/liveyoutube_counter /liveyoutube_counter
RUN apk --no-cache add tzdata
CMD ["./liveyoutube_counter"]