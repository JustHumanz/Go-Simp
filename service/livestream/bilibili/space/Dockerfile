FROM alpine
COPY --from=justhumanz/go-simp:${RELEASE_VERSION:-latest} /app/Go-Simp/bin/spacebili /spacebili
RUN apk --no-cache add tzdata
CMD ["./spacebili"]