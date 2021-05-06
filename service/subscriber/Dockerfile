FROM go-simp
FROM alpine as subscriber
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/subscriber /subscriber
CMD ["./subscriber","-Youtube","-BiliBili","-Twitter"]