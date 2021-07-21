# build web
FROM node:12-alpine as web
WORKDIR /usr/src
COPY ./web/package.json .
COPY ./web/package-lock.json .
RUN yarn install
COPY ./web .
RUN yarn build

# build server
FROM golang:1.16.5 as builder
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN rm -rf ./web
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# deploy
FROM debian:stretch
WORKDIR /cmd
COPY --from=builder /src/app ./app
COPY --from=web /usr/src/build ./web
COPY public public

ENTRYPOINT ["./app"]


