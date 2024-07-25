FROM golang:1.22-bookworm as builder

WORKDIR /build

RUN apt-get update && apt-get install --no-install-recommends -y iputils-ping

COPY . .
RUN go mod download

RUN CGO_ENABLED=0 \
    go build \
    -ldflags "-s -w" \
    -trimpath \
    -o server \
    .

FROM --platform=linux/amd64 gcr.io/distroless/static-debian12:nonroot

WORKDIR /automation_panel

COPY --from=builder --chown=nonroot:nonroot /build/server /automation_panel/server

EXPOSE 8080

CMD ["./server"]
