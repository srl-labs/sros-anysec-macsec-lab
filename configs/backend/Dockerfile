FROM golang:1.22

WORKDIR /app/server

RUN apt-get update && apt-get install --no-install-recommends -y iputils-ping

COPY . .
RUN go mod download
RUN go build

ENV HOST=0.0.0.0
EXPOSE 8080

CMD ["./backend"]
