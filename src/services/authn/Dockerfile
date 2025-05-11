FROM golang:latest AS base

WORKDIR /app

# アプリケーションのソースをコピーする
COPY .  .

FROM base AS dev

WORKDIR /app/internal

RUN go install github.com/air-verse/air@latest

CMD ["air", "-c", "air.toml"]

FROM base AS build

WORKDIR /app/internal

RUN go mod download

RUN go build -o api 

FROM scratch AS prov
COPY --from=build /app/internal/api /bin/api

EXPOSE 8080
ENTRYPOINT [ "/bin/api" ]