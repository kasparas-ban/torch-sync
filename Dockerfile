FROM golang:alpine3.18 as build
WORKDIR /server
COPY . /server
RUN go build -o /torch-server

FROM ubuntu:latest
COPY --from=build ./torch-server ./
COPY --from=build ./server/.env ./
EXPOSE 8000

RUN apt update && apt upgrade -y

ENTRYPOINT ["/torch-server"]