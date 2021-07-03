FROM alpine
COPY --from=justhumanz/go-simp /app/Go-Simp/bin/subscriber /subscriber
CMD ["./subscriber","-Youtube","-BiliBili","-Twitter","-Twitch"]