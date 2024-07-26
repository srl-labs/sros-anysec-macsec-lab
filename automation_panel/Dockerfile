FROM golang:1.22-bookworm AS be-builder

WORKDIR /build

RUN apt-get update && apt-get install --no-install-recommends -y iputils-ping

COPY backend .
RUN go mod download

RUN CGO_ENABLED=0 \
    go build \
    -ldflags "-s -w" \
    -trimpath \
    -o server \
    .

FROM node:21-alpine AS fe-builder

WORKDIR /build

COPY ./frontend .

RUN npm install
RUN npm run build

# FROM --platform=linux/amd64 gcr.io/distroless/static-debian12:nonroot

# WORKDIR /automation_panel

# COPY --from=builder --chown=nonroot:nonroot /build/server /automation_panel/server

# EXPOSE 8080

# CMD ["./server"]

FROM node:21-alpine

WORKDIR /app
COPY --from=fe-builder /build .
COPY --from=be-builder /build/server /app/server
COPY entrypoint.sh /app/entrypoint.sh

ENV HOST=0.0.0.0
EXPOSE 4173

CMD [ "/app/entrypoint.sh" ]