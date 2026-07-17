FROM golang:1.26.0 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN apt-get update && apt-get install -y \
    gcc \
    pkg-config \
    libavcodec-dev \
    libavutil-dev \
    libswscale-dev \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED=1

RUN go build -o camera-processor ./cmd

CMD ["./camera-processor"]