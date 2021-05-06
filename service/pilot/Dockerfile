FROM go-simp
FROM alpine as pilot
COPY --from=go-simp /app/Go-Simp/bin/pilot /pilot
EXPOSE 9000
CMD ["./pilot"]