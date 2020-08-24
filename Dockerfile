FROM golang:1.14 as builder

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o api .

FROM alpine:latest AS production
WORKDIR /app
COPY --from=builder /app .
CMD ["./api"]