FROM golang:1.27.1-bookworm

WORKDIR /docker

COPY . .

RUN go mod tidy

CMD ["make", "run"]