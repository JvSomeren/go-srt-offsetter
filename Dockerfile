# prepare frontend
FROM node:12.15 AS JS_BUILD
COPY webapp /webapp
WORKDIR /webapp
RUN npm install && npm run build --prod

# prepare backend
FROM golang:1.13-alpine AS GO_BUILD
RUN apk add build-base
COPY server /server
WORKDIR /server
RUN go build -o /go/bin/server

# combine and run
FROM alpine:3.11
COPY --from=JS_BUILD /webapp/public* ./webapp/
COPY --from=GO_BUILD /go/bin/server ./
VOLUME /media

ENV PORT 8080
ENV STATIC_PATH /webapp/
ENV MEDIA_PATH /media/

EXPOSE 8080
CMD ./server
