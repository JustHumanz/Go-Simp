FROM justhumanz/go-simp as go-simp
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/liveyoutube /liveyoutube
CMD ["./liveyoutube"]