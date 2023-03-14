# stage 1
FROM golang:1.18-alpine as build

WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@v1.8.5

COPY . .

RUN swag init -g internal/app/app.go -o ./gen/docs

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app/main.go

# stage 2
FROM alpine:3.12
RUN apk --update add ca-certificates tzdata

WORKDIR /app

COPY --from=build /app/gen ./gen
COPY --from=build /bin/app ./app

ENTRYPOINT [ "./app" ]
