# Arquitetura do Sistema — Bank API

## 1. Visão Geral

A arquitetura do sistema segue uma abordagem **modular em camadas**, com separação explícita de responsabilidades.

O objetivo é garantir:

* isolamento do domínio
* baixo acoplamento
* alta coesão
* facilidade de evolução

A arquitetura adotada pode ser descrita como:

> **Monólito modular com separação em camadas (Layered Architecture com influência de Clean Architecture)**

---

## 2. Estrutura Arquitetural

O sistema é dividido em quatro camadas principais:

```text
Delivery → Application → Domain
                  ↓
             Infrastructure
```

---

## 3. Descrição das Camadas

### 3.1 Domain

#### Papel

Representa o núcleo do sistema.

Contém:

* entidades
* regras de negócio
* invariantes
* contratos essenciais

#### Características

* não depende de nenhuma outra camada
* não conhece HTTP, banco ou framework
* representa o modelo financeiro do sistema

---

### 3.2 Application

#### Papel

Orquestra os casos de uso.

Responsável por:

* coordenar fluxos
* aplicar regras de processo
* garantir consistência operacional
* controlar transações

#### Características

* depende apenas do domínio
* não contém lógica de infraestrutura
* não conhece detalhes de HTTP

---

### 3.3 Infrastructure

#### Papel

Implementa detalhes técnicos.

Responsável por:

* persistência (PostgreSQL)
* acesso a dados
* integração com serviços externos (quando houver)
* geração de identificadores
* controle de tempo

#### Características

* depende do domínio
* implementa contratos definidos no domínio
* contém detalhes tecnológicos

---

### 3.4 Delivery

#### Papel

Interface externa do sistema.

Responsável por:

* exposição via HTTP
* validação de entrada
* transformação de dados (DTOs)
* mapeamento de respostas

#### Características

* depende da camada de aplicação
* não contém regras de negócio
* não acessa diretamente a infraestrutura

---

## 4. Direção das Dependências

A arquitetura respeita a seguinte regra:

> Dependências sempre apontam para dentro.

```text
Delivery → Application → Domain
Infrastructure → Domain
```

Isso garante:

* domínio isolado
* substituição de tecnologia sem impacto no core
* testabilidade

---

## 5. Organização do Código

Inicialmente, o projeto será organizado por camada:

```text
internal/
├── domain/
├── application/
├── infrastructure/
└── delivery/
```

Essa estrutura será evoluída conforme o crescimento do projeto.

---

## 6. Fluxo de Execução

Uma requisição segue o seguinte fluxo:

```text
HTTP Request
   ↓
Delivery (handler)
   ↓
Application (use case)
   ↓
Domain (regras)
   ↓
Infrastructure (persistência)
   ↓
Application
   ↓
Delivery
   ↓
HTTP Response
```

---

## 7. Controle de Transações

O controle transacional será realizado na camada de **Application**, utilizando suporte da infraestrutura.

Motivos:

* o escopo da transação é definido pelo caso de uso
* evita vazamento de lógica para camadas externas
* mantém consistência entre múltiplas operações

---

## 8. Separação de Responsabilidades

### Domain

* define o que é válido

### Application

* define o que deve acontecer

### Infrastructure

* define como acontece tecnicamente

### Delivery

* define como o sistema é acessado

---

## 9. Decisões Arquiteturais Intencionais

### 9.1 Monólito modular

O sistema não será distribuído inicialmente.

Motivo:

* reduzir complexidade
* facilitar consistência transacional
* simplificar desenvolvimento

---

### 9.2 Ausência de mensageria

Não haverá:

* filas
* eventos assíncronos
* processamento em background

Motivo:

* evitar inconsistência eventual
* manter modelo simples e determinístico

---

### 9.3 Sem CQRS inicial

Leitura e escrita compartilham o mesmo modelo.

Motivo:

* reduzir complexidade
* evitar duplicação de modelo

---

### 9.4 Sem Event Sourcing

O sistema não será baseado em eventos como fonte primária.

Motivo:

* aumento significativo de complexidade
* desnecessário para o estágio atual

---

## 10. Estratégia de Evolução

A arquitetura permite evolução gradual:

* introdução de módulos por domínio (account, transaction, etc.)
* separação mais granular por feature
* eventual introdução de mensageria
* possibilidade de extração para microserviços

---

## 11. Limitações Conhecidas

A arquitetura atual não resolve:

* escalabilidade horizontal avançada
* alta disponibilidade distribuída
* processamento assíncrono

Essas limitações são aceitáveis no estágio atual.

---

## 12. Conclusão

A arquitetura adotada prioriza:

* clareza estrutural
* consistência de domínio
* simplicidade operacional

Essa abordagem cria uma base sólida para evolução controlada, sem introduzir complexidade prematura.
