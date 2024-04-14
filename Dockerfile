FROM golang:1.22

WORKDIR /app

COPY . .

CMD ["go", "test", "-v", "./..."]