FROM golang:1.24 AS builder

WORKDIR /src
COPY . .
WORKDIR /src/admin-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/service ./cmd/admin-service

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/service /service
EXPOSE 8090
ENTRYPOINT ["/service"]
