FROM golang:1.24-bullseye as builder
WORKDIR /app
COPY . /app
EXPOSE 8080
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app ./cmd/app/main.go

FROM node:22.14-alpine as front-builder
ARG mode
RUN apk add --no-cache --virtual .gyp python3 make g++
WORKDIR /app
COPY ./frontend/package.json /frontend/package-lock.json* ./
RUN yarn install
COPY ./frontend /app
RUN yarn build --mode $mode

FROM scratch
USER 1001
WORKDIR /
COPY --from=builder /app/internal/migrations /internal/migrations
COPY --from=builder /app/app .
COPY --from=front-builder /app/dist /static
ENTRYPOINT ["./app"]