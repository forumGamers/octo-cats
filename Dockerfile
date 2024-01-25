FROM golang:1.20.1-alpine

# Set working directory ke dalam folder "bin"
WORKDIR /app/bin

RUN apk add --no-cache curl tar && \
   curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protoc-3.17.3-linux-x86_64.zip && \
   unzip protoc-3.17.3-linux-x86_64.zip -d /usr/local && \
   rm -f protoc-3.17.3-linux-x86_64.zip

# Copy file main ke dalam container
COPY ./ ./

# Build aplikasi Go
RUN go mod tidy

RUN go build main.go

# Jalankan aplikasi saat container dijalankan
CMD ["./main"]