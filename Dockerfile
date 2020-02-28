FROM node:12.15 AS JS_BUILD
COPY webapp /webapp
WORKDIR /webapp
RUN npm install && npm run build --prod

FROM golang:1.13-alphine AS GO_BUILD
RUN apk add build-base
COPY server /server
WORKDIR /server
RUN go build -o /go/bin/server

FROM alphine:3.11
COPY --from=JS_BUILD /webapp/public* ./webapp/
COPY --from=GO_BUILD /go/bin/server ./
VOLUME /media
CMD ./server
