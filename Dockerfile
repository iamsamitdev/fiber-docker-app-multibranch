# Build stage - สำหรับ development และ compile
FROM golang:1.24-alpine AS builder

# ติดตั้ง tools ที่จำเป็นสำหรับ development
RUN apk add --no-cache git

# กำหนด Working Directory ภายใน Container
WORKDIR /app

# Copy go mod files เข้าไปก่อน
# เพื่อใช้ประโยชน์จาก Docker cache layer ทำให้ไม่ต้อง download dependencies ใหม่ทุกครั้งที่แก้โค้ด
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy โค้ดทั้งหมดในโปรเจกต์เข้าไปใน container
COPY . .

# Build the application
# CGO_ENABLED=0: สำหรับ static binary
# GOOS=linux: target OS
# -a: force rebuilding
# -installsuffix cgo: เพิ่ม suffix เพื่อไม่ให้ปนกับ cache
# -o: output filename
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o fiber-app .

# Production stage - สำหรับ production deployment
FROM alpine:latest AS production

# ติดตั้ง ca-certificates สำหรับ HTTPS requests
RUN apk --no-cache add ca-certificates

# กำหนด Working Directory ภายใน Container
WORKDIR /root/

# Copy binary ที่ compiled แล้วจาก builder stage
COPY --from=builder /app/fiber-app .

# กำหนด Port ที่ Container จะทำงาน
EXPOSE 9000

# คำสั่งสำหรับรัน Fiber Application
CMD ["./fiber-app"]
