FROM node:6-alpine AS frontend-builder
ENV NODE_ENV production
WORKDIR /app
COPY ./static /app
RUN npm install
RUN npm run build

FROM golang:1.10-alpine AS backend-builder
WORKDIR /go/src/github.com/bieber/mixer/mixerserver
RUN apk add git
COPY ./mixerserver /go/src/github.com/bieber/mixer/mixerserver
RUN go get -d github.com/bieber/mixer/mixerserver
RUN mkdir /app
RUN go build -o /app/server github.com/bieber/mixer/mixerserver

FROM alpine:latest
LABEL maintainer="docker@biebersprojects.com"
EXPOSE 80

WORKDIR /app
RUN apk add ca-certificates
ENV STATIC_PATH /app/static/
COPY --from=backend-builder /app/server /app/server
COPY --from=frontend-builder /app/build /app/static
ENTRYPOINT /app/server
