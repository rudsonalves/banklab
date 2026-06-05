# Infraestrutura do Projeto — Bank API

## 1. Visão Geral

A infraestrutura do projeto tem como objetivo fornecer uma base **simples, confiável e previsível** para suportar operações financeiras.

A prioridade não é sofisticação, mas **consistência, atomicidade e rastreabilidade**.

O sistema será inicialmente:

* monolítico
* síncrono
* orientado a API REST
* com persistência relacional

---

## 2. Stack Tecnológica

### Backend

* Go (Golang)

### Banco de Dados

* PostgreSQL

### Interface de Comunicação

* HTTP/REST

---

## 3. Princípios de Infraestrutura

### 3.1 Simplicidade controlada

A infraestrutura deve evitar complexidade desnecessária:

* sem microserviços
* sem mensageria
* sem cache distribuído
* sem event sourcing

Esses elementos só serão introduzidos sob necessidade real.

---

### 3.2 Consistência forte

Operações financeiras exigem:

* transações ACID
* controle de concorrência
* integridade de dados

O PostgreSQL será utilizado como fonte de verdade.

---

### 3.3 Atomicidade de operações críticas

Operações como transferência devem ser:

* executadas dentro de transações
* protegidas contra concorrência
* livres de estados intermediários inválidos

---

### 3.4 Idempotência

Operações sensíveis devem suportar repetição segura:

* depósitos
* saques
* transferências

Isso evita duplicidade em cenários de retry.

---

### 3.5 Observabilidade mínima obrigatória

Mesmo em um sistema simples, é necessário:

* logs estruturados
* identificação de requisição (request_id)
* rastreabilidade de operações financeiras

---

## 4. Banco de Dados

### 4.1 Escolha: PostgreSQL

Motivos:

* suporte robusto a transações
* locking confiável (`SELECT FOR UPDATE`)
* integridade referencial
* maturidade operacional

---

### 4.2 Modelo inicial

O sistema será baseado em três entidades principais:

* customers
* accounts
* transactions

---

### 4.3 Estratégia de dados

* valores monetários armazenados em **centavos (BIGINT)**
* uso de **UUIDs** como identificadores
* normalização básica (sem over-engineering)

---

### 4.4 Controle de concorrência

Operações críticas utilizarão:

* `SELECT ... FOR UPDATE`
* transações explícitas

---

## 5. Gerenciamento de Migrations

As migrations serão responsáveis por:

* criação de tabelas
* alterações estruturais
* versionamento do schema

### Requisitos

* versionamento incremental
* execução automática no ambiente local
* reversibilidade (quando possível)

---

## 6. Estrutura de Inicialização

A aplicação será inicializada a partir de:

```text
cmd/api/main.go
```

Responsabilidades:

* carregar configurações
* inicializar conexão com banco
* configurar dependências
* iniciar servidor HTTP

---

## 7. Configuração

As configurações são externas ao código.

`api/.env` é a fonte de verdade para a configuração local da API. Os comandos do
Makefile carregam esse arquivo para:

* iniciar o Docker Compose
* montar a URL de migrations
* resetar o banco
* verificar readiness
* iniciar a API

Exemplos:

* APP_TOKEN
* JWT_SECRET
* JWT_ACCESS_TOKEN_DURATION
* JWT_REFRESH_TOKEN_DURATION
* TRANSACTION_PASSWORD_PEPPER
* SERVER_PORT
* DB_HOST
* DB_PORT
* DB_NAME
* DB_USER
* DB_PASSWORD

`JWT_ACCESS_TOKEN_DURATION` e `JWT_REFRESH_TOKEN_DURATION` são opcionais. Quando
omitidos, a API usa `15m` para access tokens e `168h` para sessões de refresh.
Quando configurados, ambos devem usar a sintaxe de duração do Go e ser maiores
que zero.

---

## 8. Execução Local

Ambiente local deve ser simples:

* PostgreSQL 17 via Docker
* imagem local customizada com `pg_cron`
* aplicação executada localmente

O PostgreSQL local é construído a partir de `infra/docker/postgres/Dockerfile`.
A extensão `pg_cron` é carregada por `shared_preload_libraries` e configurada
para executar jobs no banco `bank`, no timezone `America/Sao_Paulo`.

---

## 9. Segurança (nível inicial)

Neste estágio inicial:

* sem criptografia avançada
* sem gestão de segredo sofisticada
* sem autenticação complexa

O foco está no núcleo transacional.

---

## 10. Limitações intencionais

Esta infraestrutura NÃO contempla inicialmente:

* alta disponibilidade
* escalabilidade horizontal
* tolerância a falhas distribuídas
* filas assíncronas
* replicação de banco

Essas capacidades serão introduzidas conforme necessidade.

---

## 11. Evolução planejada

A infraestrutura poderá evoluir para:

* pooling avançado de conexões
* observabilidade com tracing
* cache (Redis)
* filas (RabbitMQ/Kafka)
* separação de serviços

---

## 12. Conclusão

A infraestrutura proposta prioriza:

* previsibilidade
* consistência
* simplicidade operacional

Essa abordagem reduz riscos iniciais e cria uma base sólida para evolução controlada do sistema.
