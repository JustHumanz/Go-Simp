FROM go-simp
WORKDIR /app/Go-Simp/service/fanart/
RUN go build -o fanart
CMD ["./fanart","-TwitterFanart","-BiliBiliFanart","-LewdFanart","-PixivFanArt"]