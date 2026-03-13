FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /build
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o cronex .

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /build/cronex .
CMD ["./cronex"]