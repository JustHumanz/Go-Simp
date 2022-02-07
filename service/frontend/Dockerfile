FROM alpine
COPY --from=justhumanz/go-simp:${RELEASE_VERSION:-latest} /app/Go-Simp/bin/fe /fe
RUN apk --no-cache add tzdata
CMD ["./fe"]