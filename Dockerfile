
FROM arm32v6/golang:1.15-alpine AS build-go

WORKDIR /app

COPY server/ .

RUN apk --no-cache add git gcc musl-dev

RUN go build -o server.bin

########################################################################################################################

FROM arm32v6/node:lts-alpine3.10 as build-frontend

WORKDIR /app

COPY public/ ./public
COPY src/ ./src
COPY package-lock.json package.json ./

RUN npm install @vue/cli
RUN npm i

RUN npm run build

########################################################################################################################


FROM arm32v6/alpine as certgen

WORKDIR /certs

RUN apk add openssl
RUN openssl req -new -newkey rsa:4096 -days 3650 -nodes -x509 -subj \
    "/C=US/ST=DC/L=DC/O=werabcontainers.com/CN=httpcolonslashslashwww.com" \
    -keyout ./server.key -out ./server.crt


########################################################################################################################

FROM arm32v6/nginx:alpine

COPY --from=build-frontend /app/dist/ /usr/share/nginx/html
COPY default.conf /etc/nginx/conf.d/default.conf

COPY --from=certgen /certs/server.key /etc/ssl/private/server.key
COPY --from=certgen /certs/server.crt /etc/ssl/misc/server.crt

COPY --from=build-go /app/server.bin .

WORKDIR /app

EXPOSE 443
CMD ["./server.bin"]