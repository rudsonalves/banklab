FROM golang:1.26.1-alpine AS builder

WORKDIR /src

COPY api/go.mod api/go.sum ./
RUN go mod download

COPY api/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/bank-api ./cmd/api

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S banklab \
    && adduser -S -G banklab banklab

WORKDIR /app

COPY --from=builder /out/bank-api /app/bank-api

USER banklab

EXPOSE 8080

ENTRYPOINT ["/app/bank-api"]
