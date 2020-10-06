FROM golang:1.15 as builder

WORKDIR /app
COPY ./app-credential.json ./
COPY ./src/ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -v -o subscriber

FROM alpine:latest

COPY --from=builder /app/subscriber /subscriber
COPY --from=builder /app/app-credential.json /app-credential.json
ENV GOOGLE_APPLICATION_CREDENTIALS /app-credential.json

CMD ["/subscriber"]
