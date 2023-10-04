FROM golang:1.18 as build

WORKDIR /curio-db

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/curiodb .

FROM alpine:latest

WORKDIR /app

COPY --from=build /curio-db/bin/curiodb .

EXPOSE 3141

CMD [ "./curiodb" ]