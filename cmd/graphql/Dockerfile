FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN go mod download
# Sem o build aqui, só faz run
CMD ["go", "run", "cmd/graphql/main.go"]
  