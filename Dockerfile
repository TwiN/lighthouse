FROM golang:alpine AS builder
WORKDIR /app
ADD . ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o lighthouse .
RUN apk --update add ca-certificates

FROM scratch
COPY --from=builder /app/lighthouse .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/lighthouse"]