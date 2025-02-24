FROM golang:1.23-bullseye AS builder

RUN apt-get update && apt-get install -y ffmpeg

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o server main.go

# Use a lightweight container for final execution
FROM debian:bullseye-slim

# Install only necessary packages in the final image
RUN apt-get update && apt-get clean && rm -rf /var/lib/apt/lists/* && apt-get update && apt-get install -y --no-install-recommends ffmpeg

WORKDIR /app
COPY --from=builder /app/server /app/server
COPY test /app/test
COPY ffmpeg.sh /app/ffmpeg.sh

# Ensure correct line endings and permissions
RUN sed -i 's/\r$//' /app/ffmpeg.sh && \
    chmod +x /app/ffmpeg.sh

EXPOSE 8080

CMD ["/app/server"]
