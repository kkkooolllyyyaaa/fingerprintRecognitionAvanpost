FROM golang:alpine as builder
RUN apk update --no-cache && apk add --no-cache tzdata
WORKDIR /build
COPY . .
RUN go build -o fingerprintRecognitionAvanpost/cmd/solution

FROM alpine
ENV TZ Europe/Moscow
RUN apk update --no-cache && apk add --no-cache tzdata ca-certificates
COPY --from=builder /build/app /app
# cpecify your files to test here and pass root filepath in code
COPY files/ files/
EXPOSE 8080
ENTRYPOINT ["/app"]
