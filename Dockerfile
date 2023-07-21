FROM golang:1.20.1-alpine3.17 AS build
WORKDIR /app/
COPY ./ ./
RUN ["ls", "-la", "env"]
RUN go build -o ./server ./src/cmd/http/main.go

FROM alpine:3.17 AS final
WORKDIR /app/
COPY --from=build  /app/server ./server
COPY ./env/ ./env/
RUN ["ls", "-la", "env"]
ENTRYPOINT [ "./server" ]