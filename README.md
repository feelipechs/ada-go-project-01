# orders-api

API REST para gerenciamento de pedidos com persistencia em PostgreSQL.

## Stack

- Go 1.26 + Gin
- PostgreSQL 16
- pgxpool (connection pool)
- golang-migrate (migrations)
- bcrypt (password hashing)
- Air (hot reload)

## Entidades

- Client
- Product
- Order (PENDING, PAID, CANCELED)
- OrderItem

## Endpoints

```
POST   /clients
GET    /clients
GET    /clients/:id

POST   /products
GET    /products
GET    /products/:id

POST   /orders
GET    /orders?limit=10&offset=0
GET    /orders/:id
POST   /orders/:id/pay
POST   /orders/:id/cancel

GET    /health
```

## Pre-requisitos

- Go 1.26+
- Docker (para PostgreSQL)
- golang-migrate CLI

```bash
make install-migrate
```

## Como rodar

```bash
# 1. Copiar arquivo de ambiente
cp .env.example .env

# 2. Iniciar PostgreSQL
sudo docker compose up -d

# 3. Rodar migrations
make migrate-up

# 4. Iniciar servidor
make backend
```

Ou tudo de uma vez:

```bash
sudo make up
```

Servidor disponivel em `http://localhost:8080`.

## Comandos do Makefile

| Comando           | Descricao                            |
|-------------------|--------------------------------------|
| `make up`         | Sobe banco, roda migrations, inicia API |
| `make db-up`      | Sobe container PostgreSQL            |
| `make db-down`    | Para container PostgreSQL            |
| `make migrate-up` | Executa migrations pendentes         |
| `make backend`    | Inicia servidor Go                   |
| `make down`       | Para todos os containers             |
| `make logs`       | Exibe logs do PostgreSQL             |
| `make test`       | Executa testes                       |

## Estrutura do projeto

```
.
├── main.go
├── auth/            # Hash de senhas (bcrypt)
├── config/          # Configuracoes via variaveis de ambiente
├── controllers/     # Handlers HTTP (Gin)
├── database/        # Conexao com PostgreSQL (pgxpool)
├── dto/             # Structs de request/response
├── migrations/      # Migrations SQL
├── model/           # Modelos de dominio e erros
├── repository/      # Queries SQL
├── routes/          # Registro de rotas
└── service/         # Regras de negocio
```
