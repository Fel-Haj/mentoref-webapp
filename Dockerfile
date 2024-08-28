FROM golang:alpine AS builder

WORKDIR /usr/local/app

RUN apk add --no-cache nodejs npm

COPY . .

RUN npm ci && npm run build

RUN go build -o app cmd/main.go

FROM alpine:latest

WORKDIR /usr/local/app

COPY --from=builder /usr/local/app/app .

EXPOSE 443

CMD ["./app"]
