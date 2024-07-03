# Go uygulaması için Dockerfile
FROM golang:latest

# Çalışma dizinini ayarla
WORKDIR /go/src/app   

# go.mod ve go.sum dosyalarını kopyala
COPY go.mod go.sum ./

# Modülleri indir
RUN go mod download

# Uygulama kaynak kodlarını kopyala
COPY . .

# Gereksiz bağımlılıkları temizler ve eksik bağımlılıkları ekler
RUN go mod tidy        

# Uygulamayı derle
RUN go build -o main .

EXPOSE 8040

# Uygulamayı başlat
CMD ["./main"]

