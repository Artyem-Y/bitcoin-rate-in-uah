FROM golang:1.18-alpine as builder
RUN apk add --no-cache build-base
WORKDIR /build
COPY . .
RUN go mod download
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go

FROM alpine:3
EXPOSE 8080
WORKDIR /app
COPY --from=builder /main /bin/app
RUN chmod +x /bin/app

CMD ["sh", "-c", "exec /bin/app"]  
