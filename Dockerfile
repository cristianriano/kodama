FROM golang:1.25.6-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /kodama-api ./cmd/api

FROM alpine:3.22

RUN adduser -D -H -u 10001 kodama

COPY --from=build /kodama-api /kodama-api

USER kodama

EXPOSE 8080

ENTRYPOINT ["/kodama-api"]
