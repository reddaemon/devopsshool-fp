FROM golang:latest as builder
WORKDIR /app
COPY . /app
EXPOSE 8080
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app ./cmd/app/main.go
RUN ls -la /app

FROM scratch
COPY --from=builder /app/.config .
COPY --from=builder /app/app .
ENTRYPOINT ["./app"]