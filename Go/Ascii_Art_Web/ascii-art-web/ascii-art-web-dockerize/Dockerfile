# golang üzerinden işlem yap eğer golang yoksa docker hubtan yükle 
FROM golang:latest

# app isimli bir klasör aç container orda çalışacak
WORKDIR /app

# dosyalarımı kopyalayıp oluşturduğu app klasörü içersine atıyoruz
COPY . .

# container içinde go oluşturma
RUN go build -o main .

# uygulamanın çalıştığı port
EXPOSE 8082

# dosyayı çalıştırma komutu
CMD ["./main"]