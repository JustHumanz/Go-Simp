FROM alpine
COPY --from=justhumanz/go-simp /app/Go-Simp/bin/livetwitch /livetwitch
RUN apk --no-cache add tzdata
CMD ["./livetwitch"]