FROM golang:1.18-buster
COPY . /goapp
WORKDIR /goapp

RUN go mod tidy

RUN go build -a -o /app .

FROM debian:buster-slim
WORKDIR /goapp
COPY --from=0 /app ./
RUN apt update
RUN apt -y install curl

ENTRYPOINT ["./app"]
