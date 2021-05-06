FROM go-simp
FROM alpine as spacebili
COPY --from=go-simp /app/Go-Simp/bin/spacebili /spacebili
CMD ["./spacebili"]