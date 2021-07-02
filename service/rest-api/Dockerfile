FROM justhumanz/go-simp as go-simp 
RUN apk --no-cache add tzdata

FROM alpine
COPY --from=go-simp /app/Go-Simp/bin/api /api
EXPOSE 2525
CMD ["./api"]