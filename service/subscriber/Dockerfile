ARG base_image=justhumanz/go-simp:latest
FROM ${base_image} as base

FROM alpine
COPY --from=base /app/Go-Simp/bin/subscriber /subscriber
RUN apk --no-cache add tzdata
CMD ["./subscriber","-Youtube","-BiliBili","-Twitter","-Twitch"]