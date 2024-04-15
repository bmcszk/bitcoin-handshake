FROM golang:1.22

WORKDIR /test

COPY . .

CMD ["go", "test", "-v", "./..."]
