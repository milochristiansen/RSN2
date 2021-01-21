
FROM arm32v6/golang:1.15-alpine AS build-go

WORKDIR /app

COPY server/ .

RUN apk --no-cache add git gcc musl-dev

RUN go build -o server.bin

########################################################################################################################

FROM arm32v6/node:lts-alpine3.10 as build-frontend

WORKDIR /app

COPY frontend/public/ ./public
COPY frontend/src/ ./src
COPY frontend/package-lock.json frontend/package.json ./

RUN npm install @vue/cli
RUN npm i

RUN npm run build

########################################################################################################################


FROM arm32v6/alpine as certgen

WORKDIR /certs

RUN apk add openssl
RUN openssl req -new -newkey rsa:4096 -days 3650 -nodes -x509 -subj \
    "/C=US/ST=DC/L=DC/O=httpcolonslashslashwww.com/CN=httpcolonslashslashwww.com" \
    -keyout ./server.key -out ./server.crt


########################################################################################################################

FROM arm32v6/alpine

COPY --from=build-frontend /app/dist/ /app/html

COPY --from=certgen /certs/server.key /app/cert/server.key
COPY --from=certgen /certs/server.crt /app/cert/server.crt

WORKDIR /app

COPY --from=build-go /app/server.bin .

EXPOSE 443
CMD ["/app/server.bin"]