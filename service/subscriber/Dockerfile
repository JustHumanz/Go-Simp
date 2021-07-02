FROM justhumanz/go-simp as go-simp
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/subscriber /subscriber
CMD ["./subscriber","-Youtube","-BiliBili","-Twitter","-Twitch"]