# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

ENV PORT=8080
ENV CACHE_TTL=6h

COPY --from=builder /build/calendar-go .

EXPOSE 8080

CMD ["./calendar-go"]
