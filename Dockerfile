FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=1 GOARCH=${TARGETARCH} go build -ldflags '-linkmode external -extldflags "-static"' -o mining-bot .

FROM alpine:latest

RUN apk add --no-cache wget

WORKDIR /app

COPY --from=builder /build/mining-bot /app/mining-bot

CMD ["/app/mining-bot"]