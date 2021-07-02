FROM justhumanz/go-simp as go-simp 
RUN apk --no-cache add tzdata

FROM alpine
COPY --from=go-simp /app/Go-Simp/bin/liveyoutube_counter /liveyoutube_counter
CMD ["./liveyoutube_counter"]