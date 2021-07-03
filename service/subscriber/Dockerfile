FROM alpine
COPY --from=justhumanz/go-simp /app/Go-Simp/bin/subscriber /subscriber
RUN apk --no-cache add tzdata
CMD ["./subscriber","-Youtube","-BiliBili","-Twitter","-Twitch"]