FROM golang:1.17-bullseye as builder
WORKDIR /app
COPY . /app
EXPOSE 8080
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app ./cmd/app/main.go

FROM scratch
WORKDIR /
COPY --from=builder /app/.config .
COPY --from=builder /app/app .
COPY --from=builder /app/static /static
ENTRYPOINT ["./app"]