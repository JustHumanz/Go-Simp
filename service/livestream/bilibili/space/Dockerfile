FROM go-simp
FROM alpine as spacebili
RUN apk --no-cache add tzdata
COPY --from=go-simp /app/Go-Simp/bin/spacebili /spacebili
CMD ["./spacebili"]