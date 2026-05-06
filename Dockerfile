# Stage 0
FROM golang:1.26.2-alpine3.23 AS base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


# Stage 1
FROM base AS build-server 
RUN go build -o /app/snippetbox ./cmd/web

# Stage 2
FROM scratch
WORKDIR /app

COPY --from=build-server /app/snippetbox .
COPY --from=build-server /app/ui ./ui

EXPOSE 4000

ENTRYPOINT [ "/app/snippetbox" ]
