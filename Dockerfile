FROM golang:1.24.5-alpine AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /true-markets
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY config/config.sample.yml config/config.yml
RUN GOOS=linux go build -a -tags musl -installsuffix cgo -ldflags '-extldflags "-static"' -o truemarkets cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates mailcap && addgroup -S app && adduser -S -G app --uid 1000 --home /true-markets app
USER app
COPY --chown=app:app --from=builder /true-markets/truemarkets /true-markets/
COPY --chown=app:app --from=builder /true-markets/config /true-markets/config
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
WORKDIR /true-markets
EXPOSE 8080
CMD ["./truemarkets"]
