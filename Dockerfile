FROM node:alpine as tailwind
WORKDIR /usr/local/app

COPY . .

RUN npm i
RUN npm run tailwind_prod

FROM golang:alpine

WORKDIR /usr/local/app

COPY --from=tailwind /usr/local/app/web/static/css/output.css web/static/css/output.css
COPY . .
RUN go build -v -o bin/app

CMD [ "./bin/app" ]