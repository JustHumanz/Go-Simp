FROM go-simp
FROM alpine as pilot
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/pilot /pilot
EXPOSE 9000
CMD ["./pilot"]