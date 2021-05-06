FROM go-simp
FROM alpine as subscriber
COPY --from=go-simp /app/Go-Simp/bin/subscriber /subscriber
CMD ["./subscriber","-Youtube","-BiliBili","-Twitter"]