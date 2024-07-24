FROM golang:1.22.4

WORKDIR /app/server

ENV http_proxy="http://135.245.192.7:8000"
ENV https_proxy="http://135.245.192.7:8000"
ENV no_proxy="127.0.0.1,localhost"

RUN apt-get update && apt-get install -y iputils-ping

COPY . .
RUN go mod download
RUN go build

ENV HOST=0.0.0.0
EXPOSE 8080

CMD ["./backend"]
