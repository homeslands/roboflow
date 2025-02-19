# Development with live reload
FROM golang:1.23.4-alpine AS dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD [ "air", "-c", ".air.toml" ]


# Production
FROM golang:1.23.4-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /roboflow ./cmd/server


FROM alpine:3.20 AS prod

WORKDIR /
COPY  --from=builder /roboflow /roboflow

CMD ["/roboflow"]
