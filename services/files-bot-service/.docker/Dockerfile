FROM --platform=linux/amd64 golang:1.23-alpine3.21 as build

WORKDIR /app

COPY ./ .

RUN go mod tidy

#COPY .env .env

#RUN apk --no-cache --update add build-base

RUN go build -o /app/dist/application ./app/main.go

FROM --platform=linux/amd64 alpine

USER nobody

COPY --from=build --chown=nobody:nobody /app/dist /app

WORKDIR /app

EXPOSE 50053

ENTRYPOINT ["/app/application"]