FROM golang:1.20.1-alpine3.17 AS build
ARG BUILD_DIR=/usr/src/app/
WORKDIR $BUILD_DIR
COPY ./ ./
RUN go build -o ./server ./src/cmd/http/main.go

FROM alpine:3.17
ARG BUILD_DIR=/usr/src/app/
ARG EXEC_DIR=/usr/src/app/
WORKDIR $EXEC_DIR
COPY --from=build  ${BUILD_DIR}env/* ./env/
COPY --from=build  ${BUILD_DIR}server ./server
ENTRYPOINT [ "./server" ]