FROM node:21 as build

WORKDIR /app

ENV http_proxy="http://135.245.192.7:8000"
ENV https_proxy="http://135.245.192.7:8000"
ENV no_proxy="127.0.0.1,localhost"

COPY package.json ./
COPY package-lock.json ./
RUN npm install
COPY . ./
RUN npm run build


FROM node:21-alpine

WORKDIR /app
COPY --from=build /app .

ENV HOST=0.0.0.0
EXPOSE 4173
CMD ["npm","run", "preview","--", "--host", "0.0.0.0"]
