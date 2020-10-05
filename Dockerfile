FROM golang:1.15 as builder

WORKDIR /app
COPY ./src/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o subscriber

FROM alpine:latest

COPY --from=builder /app/subscriber /subscriber
CMD ["/subscriber"]
