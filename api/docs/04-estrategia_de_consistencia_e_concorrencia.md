# Estratégia de Consistência e Concorrência — Bank API

## 1. Visão Geral

A consistência do sistema é garantida por meio de:

* transações ACID no PostgreSQL
* controle explícito de concorrência
* validações no nível da aplicação

O sistema adota um modelo de:

> **consistência forte e imediata**

Não há tolerância para inconsistência eventual em operações financeiras.

---

## 2. Objetivos

Garantir que:

* o saldo nunca seja corrompido
* operações concorrentes não gerem inconsistência
* não existam estados intermediários inválidos
* duplicidades sejam evitadas

---

## 3. Modelo de Concorrência

### 3.1 Concorrência esperada

O sistema deve suportar:

* múltiplos depósitos simultâneos
* saques concorrentes na mesma conta
* transferências simultâneas envolvendo as mesmas contas

---

### 3.2 Risco principal

Sem controle adequado:

* race conditions
* saldo negativo indevido
* duplicidade de operações

---

## 4. Estratégia Adotada

### 4.1 Transações explícitas

Todas as operações que alteram saldo devem:

* iniciar uma transação
* executar todas as operações dentro dela
* finalizar com commit ou rollback

---

### 4.2 Lock pessimista

Uso de:

```sql
SELECT ... FOR UPDATE
```

Objetivo:

* bloquear registros envolvidos
* impedir leitura/modificação concorrente

---

## 5. Operações Críticas

---

### 5.1 Depósito

#### Estratégia

1. iniciar transação
2. bloquear conta:

   ```sql
   SELECT * FROM accounts WHERE id = ? FOR UPDATE;
   ```
3. atualizar saldo
4. inserir transaction
5. commit

---

### 5.2 Saque

#### Estratégia

1. iniciar transação
2. bloquear conta
3. validar saldo
4. atualizar saldo
5. inserir transaction
6. commit

---

### 5.3 Transferência

#### Estratégia

1. iniciar transação
2. bloquear conta origem
3. bloquear conta destino
4. validar saldo origem
5. atualizar saldo origem
6. atualizar saldo destino
7. inserir duas transactions
8. commit

---

## 6. Ordem de Lock (CRÍTICO)

Para evitar deadlocks:

> As contas devem ser sempre bloqueadas na mesma ordem.

### Regra

* ordenar IDs das contas
* bloquear primeiro o menor ID
* depois o maior ID

---

## 7. Isolamento de Transação

### Nível recomendado

```text
READ COMMITTED (padrão do PostgreSQL)
```

Motivo:

* combinado com `FOR UPDATE` já garante segurança
* evita overhead desnecessário

---

### Alternativa (caso necessário)

```text
REPEATABLE READ
```

Usado apenas se surgirem inconsistências mais complexas.

---

## 8. Idempotência

### Objetivo

Evitar duplicidade em:

* retries de rede
* reenvio de requisições

---

### Estratégia

* uso de `idempotency_key`
* índice único no banco

```sql
(account_id, idempotency_key)
```

---

### Comportamento

* operação repetida → retorna sucesso anterior
* não executa novamente

---

## 9. Garantias do Sistema

### 9.1 Atomicidade

* operações completas ou inexistentes

---

### 9.2 Consistência

* regras de domínio sempre respeitadas

---

### 9.3 Isolamento

* concorrência não interfere em operações ativas

---

### 9.4 Durabilidade

* commit garante persistência

---

## 10. Falhas e Recuperação

### 10.1 Falha durante transação

* rollback automático
* nenhum estado parcial persiste

---

### 10.2 Timeout de conexão

* operação deve ser considerada não concluída
* cliente pode repetir via idempotência

---

## 11. Decisões Intencionais

### 11.1 Lock pessimista

Escolhido em vez de lock otimista.

Motivo:

* maior previsibilidade
* menor risco em sistemas financeiros

---

### 11.2 Consistência imediata

Sem uso de:

* filas
* processamento assíncrono
* eventual consistency

---

### 11.3 Banco como coordenador

O PostgreSQL é responsável por:

* serialização efetiva das operações
* controle de concorrência

---

## 12. Limitações

* não há distribuição de carga entre múltiplos bancos
* não há controle de concorrência distribuído
* escalabilidade depende do banco

---

## 13. Evolução futura

Possíveis melhorias:

* lock otimista em cenários específicos
* filas para operações não críticas
* sharding de contas
* replicação de leitura

---

## 14. Conclusão

A estratégia adotada garante:

* consistência forte
* previsibilidade
* segurança em operações financeiras

Essa abordagem privilegia correção sobre performance prematura, o que é essencial neste tipo de sistema.
