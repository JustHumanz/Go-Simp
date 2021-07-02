FROM justhumanz/go-simp as go-simp 
FROM alpine
COPY --from=go-simp /app/Go-Simp/bin/pilot /pilot
EXPOSE 9000
CMD ["./pilot"]