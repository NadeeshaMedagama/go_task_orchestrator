FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/auth-service ./cmd/auth-service

FROM alpine:3.23
WORKDIR /app
COPY --from=builder /out/auth-service /app/auth-service
EXPOSE 8081
ENTRYPOINT ["/app/auth-service"]

