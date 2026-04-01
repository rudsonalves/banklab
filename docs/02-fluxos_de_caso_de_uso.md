# Fluxos de Caso de Uso — Bank API

## 1. Visão Geral

Os casos de uso representam a execução das operações do sistema.

Cada fluxo descreve:

* entrada
* validações
* sequência de execução
* efeitos no sistema
* possíveis falhas

Todos os fluxos devem respeitar:

* invariantes do domínio
* atomicidade
* consistência de saldo

---

## 2. Criar Cliente

### Entrada

* name
* cpf
* email

### Fluxo

1. validar formato dos dados
2. verificar se CPF já existe
3. verificar se email já existe
4. criar Customer
5. persistir no banco

### Saída

* dados do cliente criado

### Erros possíveis

* CPF já cadastrado
* email já cadastrado
* dados inválidos

---

## 3. Abrir Conta

### Entrada

* customerId

### Fluxo

1. verificar se o cliente existe
2. gerar número da conta
3. gerar agência
4. criar Account com status **active**
5. saldo inicial = 0
6. persistir no banco

### Saída

* dados da conta criada

### Erros possíveis

* cliente não encontrado

---

## 4. Consultar Saldo

### Entrada

* accountId

### Fluxo

1. verificar se a conta existe
2. retornar saldo atual

### Saída

* balance

### Erros possíveis

* conta não encontrada

---

## 5. Depósito

### Entrada

* accountId
* amount
* description (opcional)
* idempotencyKey (opcional)

### Fluxo

1. validar valor > 0
2. verificar idempotência (se aplicável)
3. verificar se conta existe
4. verificar se conta está ativa
5. iniciar transação
6. atualizar saldo (saldo + amount)
7. registrar Transaction (type: deposit)
8. persistir alterações
9. finalizar transação

### Saída

* saldo atualizado

### Erros possíveis

* valor inválido
* conta não encontrada
* conta inativa
* duplicidade (idempotência)

---

## 6. Saque

### Entrada

* accountId
* amount
* description (opcional)
* idempotencyKey (opcional)

### Fluxo

1. validar valor > 0
2. verificar idempotência (se aplicável)
3. verificar se conta existe
4. verificar se conta está ativa
5. iniciar transação
6. verificar saldo suficiente
7. atualizar saldo (saldo - amount)
8. registrar Transaction (type: withdraw)
9. persistir alterações
10. finalizar transação

### Saída

* saldo atualizado

### Erros possíveis

* valor inválido
* saldo insuficiente
* conta não encontrada
* conta inativa
* duplicidade (idempotência)

---

## 7. Transferência

### Entrada

* fromAccountId
* toAccountId
* amount
* description (opcional)
* idempotencyKey (opcional)

### Fluxo

1. validar valor > 0
2. validar contas diferentes
3. verificar idempotência (se aplicável)
4. verificar se contas existem
5. verificar se ambas estão ativas
6. iniciar transação
7. bloquear ambas as contas para atualização
8. verificar saldo suficiente na origem
9. debitar conta origem
10. creditar conta destino
11. registrar Transaction:

    * transfer_out (origem)
    * transfer_in (destino)
12. persistir alterações
13. finalizar transação

### Saída

* saldo da conta origem atualizado

### Erros possíveis

* valor inválido
* contas iguais
* saldo insuficiente
* conta não encontrada
* conta inativa
* duplicidade (idempotência)

---

## 8. Consultar Extrato

### Entrada

* accountId
* limit (opcional)
* offset (opcional)

### Fluxo

1. verificar se conta existe
2. buscar transações da conta
3. ordenar por data (desc)
4. aplicar paginação

### Saída

* lista de transações

### Erros possíveis

* conta não encontrada

---

## 9. Garantias dos Fluxos

Todos os fluxos que alteram saldo devem garantir:

### 9.1 Atomicidade

* operação completa ou não ocorre

### 9.2 Consistência

* saldo nunca entra em estado inválido

### 9.3 Isolamento

* concorrência não corrompe dados

### 9.4 Durabilidade

* uma vez concluída, a operação persiste

---

## 10. Padrões Comuns

### 10.1 Validação antecipada

Sempre antes de iniciar transação:

* validar entrada
* validar estado da conta

---

### 10.2 Transação explícita

Toda operação financeira:

* inicia transação
* executa alterações
* finaliza com commit ou rollback

---

### 10.3 Registro obrigatório

Toda alteração de saldo:

* gera uma Transaction correspondente

---

### 10.4 Idempotência

Quando presente:

* impede duplicidade
* garante segurança em retries

---

## 11. Conclusão

Os fluxos definidos garantem:

* comportamento previsível
* consistência operacional
* aderência ao modelo de domínio

Esse documento conecta o domínio à execução prática do sistema.
