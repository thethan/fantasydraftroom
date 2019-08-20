FROM golang:1.12.7-stretch


WORKDIR /app

COPY . .

RUN go build -o main cmd/server/main.go

ADD .env /home/forge/fantasydraftroom.com/

EXPOSE 8081

CMD ["./main"]