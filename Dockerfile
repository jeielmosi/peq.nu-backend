FROM golang:1.20.1-alpine3.17 AS build
WORKDIR /app
COPY ./ ./
RUN go build -o ./server ./src/cmd/http/main.go

FROM alpine:3.17
WORKDIR /app
COPY --from=build /app/envs ./envs
COPY --from=build /app/server ./server
ENTRYPOINT [ "./server" ]