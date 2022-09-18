FROM node:6-alpine AS frontend-builder
WORKDIR /app
COPY ./static /app
RUN npm install
ENV NODE_ENV production
RUN npm run build

FROM golang:1.18-alpine AS backend-builder
RUN mkdir /app
WORKDIR /app
RUN apk add git
COPY ./mixerserver ./
RUN go build -o server

FROM alpine:latest
LABEL maintainer="docker@biebersprojects.com"
EXPOSE 80

WORKDIR /app
RUN apk add ca-certificates
ENV STATIC_PATH /app/static/
COPY --from=backend-builder /app/server /app/server
COPY --from=frontend-builder /app/build /app/static
ENTRYPOINT /app/server
