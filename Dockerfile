FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /roboflow ./cmd

FROM alpine:3.20

COPY  --from=builder /roboflow /roboflow

WORKDIR /

CMD ["/roboflow"]
