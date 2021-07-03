FROM alpine
COPY --from=justhumanz/go-simp /app/Go-Simp/bin/fe /fe
RUN apk --no-cache add tzdata
CMD ["./fe"]