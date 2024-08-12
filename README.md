# Product Microservice

Este projeto é um microserviço de gerenciamento de produtos desenvolvido em Go, seguindo os princípios de Domain-Driven
Design (DDD) e Clean Architecture, utilizando PostgreSQL como banco de dados.

## Descrição

O Product Microservice é responsável por gerenciar as operações relacionadas a produtos em um sistema distribuído. Ele
fornece uma API RESTful para criar, ler, atualizar e deletar informações de produtos.

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

## Executando o Projeto

Para executar este projeto, siga os passos abaixo:

1. Certifique-se de ter o Docker e o Docker Compose instalados em sua máquina.

2. Clone o repositório:

   `git clone https://github.com/HaroldoFV/product-service`

   `cd product-service`

3. Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:
   ```  
   DB_DRIVER=postgres
   DB_HOST=postgres
   DB_PORT=5432
   DB_USER=seu_usuario
   DB_PASSWORD=sua_senha
   DB_NAME=nome_do_banco
   WEB_SERVER_PORT=8000


4. Inicie os serviços usando Docker Compose:

   `docker-compose up -d`

   Isso irá iniciar o banco de dados PostgreSQL, executar as migrações e iniciar a aplicação.

5. A aplicação estará disponível em `http://localhost:8000/docs/index.html`.

## Executando os Testes

Para executar os testes, siga estas etapas:

1. Certifique-se de que o container de teste do PostgreSQL(product_db_test) está em execução:

   `docker-compose ps`

2. Execute os testes usando o seguinte comando:

   `go test ./... -v`

Este comando executará todos os testes no projeto, incluindo testes de unidade e integração.

Nota: Os testes de integração usarão o banco de dados de teste (postgres_test) que está configurado para rodar na porta
5433.


## Diagramas
<img width="1089" alt="Screenshot 2024-08-11 at 11 49 03 PM" src="https://github.com/user-attachments/assets/a038e810-bc0d-4822-9cab-ef77449f8bba">



