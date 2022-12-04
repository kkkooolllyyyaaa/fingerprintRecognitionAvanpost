FROM golang:alpine as builder
RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /build
COPY . .
RUN go build -o fingerprintRecognitionAvanpost/cmd/test

FROM alpine
ENV TZ Europe/Moscow
RUN apk update --no-cache && apk add --no-cache tzdata ca-certificates
COPY --from=builder /build/app /app
COPY configs/ configs/
EXPOSE 8080
ENTRYPOINT ["/app"]
