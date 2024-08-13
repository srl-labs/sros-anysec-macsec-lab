FROM golang:1.21-alpine AS be-builder

WORKDIR /build

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

# Resulting container image
FROM node:21-alpine

WORKDIR /app
COPY --from=fe-builder /build .
COPY --from=be-builder /build/server /app/server
COPY entrypoint.sh /app/entrypoint.sh

ENV HOST=0.0.0.0
EXPOSE 4173

CMD [ "/app/entrypoint.sh" ]