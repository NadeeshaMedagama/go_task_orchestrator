FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/api-gateway ./cmd/api-gateway

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/api-gateway /app/api-gateway
EXPOSE 8080
ENTRYPOINT ["/app/api-gateway"]

