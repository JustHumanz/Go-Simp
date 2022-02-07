FROM alpine
COPY --from=justhumanz/go-simp:${RELEASE_VERSION:-latest} /app/Go-Simp/bin/pilot /pilot
EXPOSE 9000
CMD ["./pilot"]