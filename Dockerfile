FROM golang:1.22-alpine

LABEL maintainer="wooseop.kim.dev@gmail.com"

WORKDIR /var/app

RUN apk add make

COPY go.* ./
RUN go mod download

COPY . .

CMD ["make", "run"]
