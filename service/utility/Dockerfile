FROM go-simp
FROM alpine as utility
COPY --from=go-simp /app/Go-Simp/bin/utility /utility
CMD ["./utility"]