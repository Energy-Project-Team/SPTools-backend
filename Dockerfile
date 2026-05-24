# Copyright (C) 2026 Energy Project Team
# SPDX-License-Identifier: AGPL-3.0-only
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/app .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/app ./app
COPY --from=builder /app/spworlds.yaml ./spworlds.yaml

EXPOSE 4042

CMD ["./app"]