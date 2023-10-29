FROM golang:1.18 as build

WORKDIR /curiodb

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/curiodb .

FROM alpine:latest

WORKDIR /app

COPY --from=build /curiodb/bin/curiodb .

EXPOSE 3141

CMD [ "./curiodb" ]