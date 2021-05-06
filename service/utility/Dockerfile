FROM go-simp
FROM alpine as utility
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/utility /utility
CMD ["./utility"]