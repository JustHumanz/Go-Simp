FROM alpine
COPY --from=justhumanz/go-simp /app/Go-Simp/bin/spacebili /spacebili
RUN apk --no-cache add tzdata
CMD ["./spacebili"]