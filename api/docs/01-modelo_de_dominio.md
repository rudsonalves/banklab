# Modelo de Domínio — Bank API

## 1. Visão Geral

O domínio do sistema é baseado no conceito de:

> **controle de saldo por meio de registros de movimentações financeiras**

Neste modelo:

* a **Transaction** é o evento fundamental
* a **Account** representa o estado
* o **saldo é derivado de operações válidas**

---

## 2. Entidades do Domínio

### 2.1 Customer

Representa o titular do sistema.

#### Atributos

* id
* name
* cpf
* email
* createdAt

#### Regras

* CPF deve ser único
* email deve ser único

---

### 2.2 Account

Representa uma conta bancária.

#### Atributos

* id
* customerId
* number
* branch
* balance (em centavos)
* status
* createdAt

#### Status possíveis

* active
* inactive
* blocked

---

### 2.3 Transaction

Representa qualquer movimentação financeira.

#### Atributos

* id
* accountId
* type
* amount (em centavos)
* description
* relatedAccountId (opcional)
* idempotencyKey (opcional)
* createdAt

#### Tipos possíveis

* deposit
* withdraw
* transfer_in
* transfer_out

---

## 3. Invariantes do Domínio

### 3.1 Valor monetário

* deve ser inteiro (centavos)
* deve ser maior que zero

---

### 3.2 Conta

* deve existir para qualquer operação
* deve estar com status **active** para movimentações

---

### 3.3 Saldo

* nunca pode se tornar inconsistente
* não pode ser alterado sem transação correspondente

---

### 3.4 Transação

* toda transação deve estar associada a uma conta
* não existe transação sem impacto financeiro

---

### 3.5 Transferência

* conta de origem ≠ conta de destino
* deve gerar duas transações:

  * débito na origem
  * crédito no destino

---

## 4. Regras de Negócio

### 4.1 Depósito

* valor deve ser positivo
* saldo da conta é incrementado
* uma transação do tipo **deposit** é registrada

---

### 4.2 Saque

* valor deve ser positivo
* saldo deve ser suficiente
* saldo da conta é decrementado
* uma transação do tipo **withdraw** é registrada

---

### 4.3 Transferência

* valor deve ser positivo
* saldo da conta de origem deve ser suficiente
* operação deve ser atômica
* duas transações devem ser registradas:

  * **transfer_out** (origem)
  * **transfer_in** (destino)

---

### 4.4 Extrato

* deve refletir todas as transações da conta
* ordenado por data

---

## 5. Consistência de Saldo

O saldo da conta deve obedecer:

```text
saldo = soma(deposits + transfer_in) - soma(withdraw + transfer_out)
```

---

## 6. Idempotência

Operações críticas podem possuir `idempotencyKey`.

Regras:

* mesma chave não pode gerar duplicidade
* a operação deve ser considerada já executada

Aplicável a:

* depósito
* saque
* transferência

---

## 7. Comportamento sob Concorrência

* operações que alteram saldo devem ser protegidas
* duas operações simultâneas não podem corromper o saldo
* o sistema deve garantir consistência final

---

## 8. Regras de Integridade

* nenhuma alteração de saldo ocorre sem registro em Transaction
* não existem estados intermediários inválidos
* operações devem ser completas ou não ocorrer

---

## 9. Decisões Intencionais

### 9.1 Saldo armazenado

O saldo é armazenado diretamente na conta por eficiência.

Trade-off:

* ganho de performance
* necessidade de controle rigoroso de consistência

---

### 9.2 Transação como fonte de auditoria

Transactions são:

* históricas
* imutáveis
* auditáveis

---

### 9.3 Modelo simplificado

O sistema não contempla:

* reversão automática de transações
* estornos complexos
* múltiplas moedas

---

## 10. Limitações do Modelo

* não há suporte a saldo negativo (descoberto)
* não há controle de limites transacionais
* não há distinção entre tipos de conta (PF/PJ)

---

## 11. Conclusão

O modelo de domínio define um sistema:

* centrado em transações
* consistente por construção
* simples, porém correto

Essa base é suficiente para suportar evolução sem comprometer integridade financeira.
