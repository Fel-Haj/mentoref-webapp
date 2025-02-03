FROM --platform=$BUILDPLATFORM golang:alpine AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

RUN apk add --no-cache nodejs npm

COPY . .

RUN npm ci && npm run build

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o app cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=build /app .

EXPOSE 443

CMD ["./app"]
