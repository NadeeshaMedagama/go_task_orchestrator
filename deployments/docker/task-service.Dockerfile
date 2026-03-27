FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /out/task-service ./cmd/task-service

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/task-service /app/task-service
EXPOSE 8082
ENTRYPOINT ["/app/task-service"]

