# Sử dụng golang dựa trên alpine cho kích thước nhỏ gọn
FROM golang:alpine

# Thiết lập thư mục làm việc
WORKDIR /app

# Sao chép go mod và sum files
COPY go.mod go.sum ./

# Tải tất cả các phụ thuộc
RUN go mod download

# Sao chép mã nguồn từ thư mục hiện tại vào Container
COPY . .

# Biên dịch ứng dụng
RUN go build -o main .

# Chạy ứng dụng
CMD ["/app/main"]
