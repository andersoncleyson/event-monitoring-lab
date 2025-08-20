FROM golang:1.19-alpine as builder

WORKDIR /app
COPY . .

RUN go mod init exporter && go mod tidy
RUN go build -o /exporter-bin .

FROM alpine:latest
COPY --from=builder /exporter-bin /exporter-bin

EXPOSE 8081

CMD ["/exporter-bin"]