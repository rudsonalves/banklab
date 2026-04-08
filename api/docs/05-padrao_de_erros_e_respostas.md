# Padrão de Erros e Respostas — Bank API

## 1. Visão Geral

Este documento define o padrão de:

* respostas de sucesso
* respostas de erro
* estrutura de mensagens
* mapeamento para HTTP status codes

O objetivo é garantir:

* consistência entre endpoints
* previsibilidade para clientes
* facilidade de tratamento de erros

---

## 2. Princípios

### 2.1 Simplicidade

* estrutura única para respostas
* sem campos desnecessários
* fácil de consumir

---

### 2.2 Consistência

* mesmo formato para todos os endpoints
* mesmos códigos para mesmos tipos de erro

---

### 2.3 Clareza

* mensagens objetivas
* códigos de erro explícitos

---

## 3. Estrutura de Resposta

### 3.1 Resposta de Sucesso

```json
{
  "data": { ... },
  "error": null
}
```

---

### 3.2 Resposta de Erro

```json
{
  "data": null,
  "error": {
    "code": "ERROR_CODE",
    "message": "Descrição do erro",
    "details": {}
  }
}
```

---

## 4. Campos

### 4.1 data

* contém o resultado da operação
* pode ser objeto ou lista

---

### 4.2 error.code

Identificador estável do erro.

Exemplos:

* CUSTOMER_ALREADY_EXISTS
* ACCOUNT_NOT_FOUND
* INSUFFICIENT_FUNDS
* INVALID_AMOUNT

---

### 4.3 error.message

Mensagem legível para humanos.

* não precisa ser traduzida neste momento
* não deve conter detalhes sensíveis

---

### 4.4 error.details (opcional)

Informações adicionais estruturadas.

Exemplo:

```json
{
  "field": "amount",
  "reason": "must be greater than zero"
}
```

---

## 5. HTTP Status Codes

### 5.1 Sucesso

| Código | Uso                   |
| ------ | --------------------- |
| 200    | operação bem-sucedida |
| 201    | recurso criado        |

---

### 5.2 Erros de cliente

| Código | Uso                        |
| ------ | -------------------------- |
| 400    | requisição inválida        |
| 404    | recurso não encontrado     |
| 409    | conflito (ex: duplicidade) |
| 422    | regra de negócio violada   |

---

### 5.3 Erros de servidor

| Código | Uso          |
| ------ | ------------ |
| 500    | erro interno |

---

## 6. Mapeamento de Erros de Domínio

| Erro de domínio      | HTTP | Code                    |
| -------------------- | ---- | ----------------------- |
| cliente já existe    | 409  | CUSTOMER_ALREADY_EXISTS |
| conta não encontrada | 404  | ACCOUNT_NOT_FOUND       |
| saldo insuficiente   | 422  | INSUFFICIENT_FUNDS      |
| valor inválido       | 400  | INVALID_AMOUNT          |
| conta inativa        | 422  | ACCOUNT_INACTIVE        |
| operação duplicada   | 409  | DUPLICATE_REQUEST       |

---

## 7. Regras de Uso

### 7.1 Nunca misturar sucesso e erro

* se houver erro → `data = null`
* se houver sucesso → `error = null`

---

### 7.2 Não expor detalhes internos

Evitar:

* stack trace
* SQL
* detalhes de infraestrutura

---

### 7.3 Padronização de códigos

* códigos em UPPER_SNAKE_CASE
* estáveis ao longo do tempo

---

### 7.4 Mensagens não são contrato

* clientes devem usar `error.code`
* não depender do texto

---

## 8. Exemplos

---

### 8.1 Sucesso — Depósito

```json
{
  "data": {
    "account_id": "uuid",
    "balance": 150000
  },
  "error": null
}
```

---

### 8.2 Erro — Saldo insuficiente

```json
{
  "data": null,
  "error": {
    "code": "INSUFFICIENT_FUNDS",
    "message": "Insufficient balance",
    "details": {}
  }
}
```

---

### 8.3 Erro — Campo inválido

```json
{
  "data": null,
  "error": {
    "code": "INVALID_AMOUNT",
    "message": "Amount must be greater than zero",
    "details": {
      "field": "amount"
    }
  }
}
```

---

## 9. Decisões Intencionais

### 9.1 Estrutura única

Evita múltiplos formatos de resposta.

---

### 9.2 Código separado da mensagem

Permite:

* internacionalização futura
* estabilidade para clientes

---

### 9.3 Uso de 422

Utilizado para regras de negócio:

* saldo insuficiente
* conta inativa

---

## 10. Limitações

* não há versionamento de erro
* não há multilíngue
* não há correlação com logs (ainda)

---

## 11. Evolução futura

Possíveis melhorias:

* adicionar `request_id`
* padronizar erro por domínio
* suporte a múltiplos idiomas
* catálogo de erros centralizado

---

## 12. Conclusão

O padrão definido:

* é simples
* é consistente
* é suficiente para evolução

Ele garante previsibilidade sem adicionar complexidade desnecessária.
