Olá devs!
Agora é a hora de botar a mão na massa. Para este desafio, você precisará criar o usecase de listagem das orders.
Esta listagem precisa ser feita com:
- Endpoint REST (GET /order)
- Service ListOrders com GRPC
- Query ListOrders GraphQL
Não esqueça de criar as migrações necessárias e o arquivo api.http com a request para criar e listar as orders.

Para a criação do banco de dados, utilize o Docker (Dockerfile / docker-compose.yaml), com isso ao rodar o comando docker compose up tudo deverá subir, preparando o banco de dados.
Inclua um README.md com os passos a serem executados no desafio e a porta em que a aplicação deverá responder em cada serviço.


# 🚀 Full Cycle Clean Architecture

![Go Version](https://img.shields.io/badge/Go-1.24.3-blue)
![Docker](https://img.shields.io/badge/docker-ready-brightgreen)
![GraphQL](https://img.shields.io/badge/graphql-enabled-purple)
![gRPC](https://img.shields.io/badge/gRPC-enabled-blue)
![REST](https://img.shields.io/badge/REST-enabled-green)

> Projeto com arquitetura limpa, multi-camada e suporte a REST, gRPC e GraphQL.

---

## Índice

* [Sobre o Projeto](#sobre-o-projeto)
* [Pré-requisitos](#pré-requisitos)
* [Instalação](#instalação)
* [Gerando código gRPC](#gerando-código-grpc)
* [Executando os Serviços](#executando-os-serviços)
* [Acessando os Endpoints](#acessando-os-endpoints)
* [Estrutura do Projeto](#estrutura-do-projeto)
* [Principais Dependências](#principais-dependências)
* [Contribuindo](#contribuindo)
* [FAQ](#faq)

---

## Sobre o Projeto

Esse repositório demonstra como estruturar uma aplicação Go utilizando **Clean Architecture**. Inclui API REST, server gRPC e GraphQL, com banco de dados PostgreSQL, migrations e integração total via Docker Compose.
Ideal para estudo de arquiteturas modernas, integração de múltiplos protocolos e testes de interoperabilidade.

---

## Pré-requisitos

* Go **>= 1.24**
* [Docker](https://www.docker.com/) + [Docker Compose](https://docs.docker.com/compose/)
* [Protocol Buffers (protoc)](https://github.com/protocolbuffers/protobuf)
* Ferramentas gRPC para Go:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

---

## Instalação

```
git clone https://github.com/JeanGrijp/Full-Cycle-Desafio-Clean-Architecture.git
cd Full-Cycle-Desafio-Clean-Architecture
go mod tidy
```

---

## Gerando código gRPC

Sempre que alterar o arquivo `proto/order.proto`, gere o código Go do gRPC:

```
protoc --go_out=internal/infra/grpc/pb \
  --go-grpc_out=internal/infra/grpc/pb \
  --proto_path=proto \
  proto/order.proto
```

---

## Executando os Serviços

Suba tudo com um comando só:

```
docker-compose up --build
```

---

## Acessando os Endpoints

* **REST:** [`http://localhost:8080`](http://localhost:8080)
* **GraphQL Playground:** [`http://localhost:8081`](http://localhost:8081)
* **gRPC:** `localhost:50051` (use BloomRPC, grpcurl ou outro client gRPC)

---

### Testando a API REST

Você pode testar o endpoint REST usando o arquivo [`api.http`](./api.http) (recomendado para VS Code ou plugins de HTTP client) ou via curl:

```http
GET http://localhost:8080/orders
Accept: application/json
```
ou com o curl:
```
curl -X GET http://localhost:8080/orders -H "Accept: application/json"
```

### Testando o graphQL
Acesse o GraphQL Playground em [http://localhost:8081/](http://localhost:8081/) e execute a query:

```graphql
query {
  listOrders {
    id
    customerName
    amount
    status
    createdAt
  }
}
````


### Testando o gRPC
Use o BloomRPC ou grpcurl para testar o serviço gRPC. Exemplo com grpcurl:

```sh
grpcurl -plaintext -d '{}' localhost:50051 orderpb.OrderService/ListOrders
```

---

## Estrutura do Projeto

```
.
├── cmd
│   ├── graphql    # Main GraphQL
│   ├── grpc       # Main gRPC
│   └── rest       # Main REST
├── internal
│   ├── infra
│   │   └── grpc
│   │       └── pb   # Arquivos .pb.go gerados
│   ├── repository   # Repositórios de dados
│   └── usecase      # Casos de uso
├── proto           # Arquivo .proto do gRPC
├── migrations      # Scripts SQL para o banco
└── ...
```

---

## Principais Dependências

* [Go >= 1.24.3](https://go.dev/dl/)
* [Docker](https://www.docker.com/)
* [Protobuf](https://github.com/protocolbuffers/protobuf)
* [gqlgen](https://gqlgen.com/)
* [chi](https://github.com/go-chi/chi)
* [PostgreSQL](https://www.postgresql.org/)

---

## Contribuindo

Quer contribuir? Fique à vontade!
Siga os passos abaixo:

1. **Fork** este repositório
2. Crie uma branch para sua feature ou correção:

   ```sh
   git checkout -b feature/nome-da-sua-feature
   ```
3. Commit suas mudanças:

   ```sh
   git commit -m 'feat: descreva sua feature'
   ```
4. Push para sua branch:

   ```sh
   git push origin feature/nome-da-sua-feature
   ```
5. Abra um **Pull Request**

Não esqueça de rodar os testes e garantir que tudo está funcionando antes de abrir o PR!

---

## FAQ

* **O projeto funciona em Windows, Linux e macOS?**
  Sim! Apenas adapte o caminho do `protoc` se necessário e garanta que o Docker está rodando normalmente.

* **Alterei o .proto, o que fazer?**
  Rode o comando de geração do gRPC descrito acima.

* **O banco não sobe ou containers não conectam?**
  Tente rodar:

  ```sh
  docker-compose down -v
  docker-compose up --build
  ```

  Isso apaga os volumes e reinicializa tudo do zero.

* **Preciso de Go instalado para rodar o projeto?**
  Só para desenvolvimento local. No Docker, tudo já roda sem precisar instalar Go.

---

## Licença

Esse projeto está sob a licença MIT.

---

Curtiu? Deixe sua ⭐️ no projeto ou contribua para ajudar a comunidade!
Dúvidas, sugestões ou bugs, só abrir uma issue ou chamar! 😃
