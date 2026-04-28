FROM golang:1.24 AS builder

WORKDIR /src
COPY . .
WORKDIR /src/merchant-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/service ./cmd/merchant-service

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/service /service
EXPOSE 8083
ENTRYPOINT ["/service"]
