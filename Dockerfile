FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git (diperlukan oleh beberapa modul Go)
RUN apk add --no-cache git

# Download dependensi lebih dulu (manfaatkan layer cache Docker)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary — static linking untuk binary yang portable
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o main .

# ==========================================
# Stage 2: Runtime image minimalis
# ==========================================
FROM alpine:3.19

# ca-certificates untuk HTTPS (koneksi ke Render Postgres), tzdata untuk timezone
RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Jakarta

WORKDIR /root/

COPY --from=builder /app/main .

# Render akan secara otomatis menginjeksi $PORT — EXPOSE ini hanya dokumentasi
EXPOSE 8080

CMD ["./main"]
