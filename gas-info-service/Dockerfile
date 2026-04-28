FROM golang:1.24 AS builder

WORKDIR /src
COPY . .
WORKDIR /src/gas-info-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/service ./cmd/gas-info-service

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/service /service
EXPOSE 8082
ENTRYPOINT ["/service"]
