# Contornos do Sistema — Bank API

## 1. Objetivo do Sistema

O sistema tem como objetivo implementar um núcleo bancário simplificado, capaz de:

* gerenciar clientes
* manter contas bancárias
* registrar movimentações financeiras
* garantir consistência de saldo

O foco está no **controle transacional confiável**, não em funcionalidades periféricas.

---

## 2. Natureza do Sistema

Este sistema é, essencialmente:

> Um sistema de controle de saldo baseado em registros de movimentações financeiras.

Ou seja:

* o saldo é consequência
* a transação é o elemento central

---

## 3. Escopo Funcional

### 3.1 Incluído no escopo

O sistema será responsável por:

#### Clientes

* criação de cliente
* identificação por CPF e email

#### Contas

* abertura de conta
* consulta de saldo
* controle de status da conta

#### Movimentações

* depósito
* saque
* transferência entre contas
* registro de todas as operações

#### Extrato

* listagem de transações por conta

---

### 3.2 Fora do escopo (neste momento)

O sistema NÃO será responsável por:

* autenticação e login
* integração com sistemas externos (Pix, TED, etc.)
* antifraude
* análise de risco
* notificações (email, push)
* gestão de múltiplas moedas
* reconciliação bancária
* liquidação externa

---

## 4. Entidades Centrais

O sistema será estruturado em torno de três entidades principais:

### 4.1 Customer

Representa o titular.

### 4.2 Account

Representa a conta bancária.

### 4.3 Transaction

Representa qualquer movimentação financeira.

---

## 5. Responsabilidades do Sistema

O sistema deve garantir:

### 5.1 Integridade financeira

* não permitir inconsistência de saldo
* garantir que toda movimentação seja registrada

---

### 5.2 Atomicidade

* operações críticas devem ser indivisíveis
* especialmente transferências

---

### 5.3 Rastreabilidade

* todas as operações devem ser auditáveis
* nenhuma alteração de saldo sem registro

---

### 5.4 Consistência

* o estado do sistema deve ser sempre válido
* mesmo sob concorrência

---

## 6. Modelo Operacional

### 6.1 Operações síncronas

Todas as operações são:

* executadas de forma síncrona
* concluídas no momento da requisição

---

### 6.2 Fonte de verdade

O banco de dados relacional é:

* a única fonte de verdade
* responsável pela consistência final

---

### 6.3 Ausência de eventual consistency

Neste estágio:

* não há processamento assíncrono
* não há inconsistência temporária aceitável

---

## 7. Limites do Sistema

### 7.1 Limite técnico

O sistema é um monólito:

* uma única aplicação
* um único banco de dados

---

### 7.2 Limite de responsabilidade

O sistema controla apenas:

* estado interno das contas
* movimentações internas

Não controla:

* dinheiro real externo
* liquidação bancária

---

## 8. Regras Fundamentais

### 8.1 Toda movimentação gera transação

Não existe alteração de saldo sem registro.

---

### 8.2 Transferência é composta

Uma transferência sempre gera:

* débito na origem
* crédito no destino

---

### 8.3 Valores são inteiros

* representados em centavos
* sem uso de ponto flutuante

---

### 8.4 Operações devem ser idempotentes

Repetições não devem causar duplicidade.

---

## 9. Estratégia de Evolução

O sistema será evoluído de forma incremental:

1. estabilizar o núcleo transacional
2. garantir consistência e testes
3. expandir funcionalidades gradualmente

---

## 10. Decisões Arquiteturais Intencionais

* evitar abstrações prematuras
* priorizar clareza sobre generalização
* crescer a estrutura sob demanda
* manter domínio como centro do sistema

---

## 11. Conclusão

Os contornos definidos estabelecem um sistema:

* pequeno, porém correto
* simples, porém consistente
* limitado, porém bem delimitado

Essa base é suficiente para evoluir o sistema sem comprometer sua integridade estrutural.
