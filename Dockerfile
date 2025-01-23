FROM golang:1.23.5-alpine3.20 as builder

ENV APP_HOME=/go/src/web

WORKDIR "${APP_HOME}"

COPY ./go.mod ./go.sum ./

RUN go mod download
RUN go mod verify

COPY ./internal ./internal
COPY ./cmd ./cmd
COPY ./docs ./docs

RUN go build -o ./bin/web ./cmd

FROM alpine:latest

ENV APP_HOME=/go/src/web
RUN mkdir -p "${APP_HOME}"

WORKDIR "${APP_HOME}"

COPY --from=builder "${APP_HOME}"/bin/web "${APP_HOME}"

ENV PORT=8080

EXPOSE ${PORT}

CMD ["./web"]
