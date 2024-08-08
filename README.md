# Product Microservice

Este projeto é um microserviço de gerenciamento de produtos desenvolvido em Go, seguindo os princípios de Domain-Driven Design (DDD) e Clean Architecture, utilizando PostgreSQL como banco de dados.

## Descrição

O Product Microservice é responsável por gerenciar as operações relacionadas a produtos em um sistema distribuído. Ele fornece uma API RESTful para criar, ler, atualizar e deletar informações de produtos.

## Tecnologias Utilizadas

- Go (Golang)
- PostgreSQL
- Docker

## Arquitetura

Este projeto segue os princípios de Clean Architecture e Domain-Driven Design:

- **Domain Layer**: Contém as entidades de negócio e regras de domínio.
- **Use Case Layer**: Implementa a lógica de aplicação e casos de uso.
- **Interface Layer**: Gerencia a comunicação com o mundo externo (API, banco de dados).
- **Infrastructure Layer**: Fornece implementações concretas para interfaces definidas em camadas superiores.
 