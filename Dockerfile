# syntax=docker/dockerfile:1

FROM golang:1.19-alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY client/ ./client/
RUN go build -o /websocket-chat-server

EXPOSE 8080

CMD [ "/websocket-chat-server" ]