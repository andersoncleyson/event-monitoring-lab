FROM golang:1.23-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /exporter-bin .

FROM alpine:latest

COPY --from=builder /exporter-bin /exporter-bin

EXPOSE 8081

CMD ["/exporter-bin"]