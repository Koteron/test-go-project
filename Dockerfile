FROM golang:1.24
COPY . /app
EXPOSE 8080
WORKDIR /app
CMD ["go", "run", "./cmd/main.go"]