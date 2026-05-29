# Backlog: senha transacional

## 1. Objetivo

Implementar a credencial de senha transacional como primeiro fator de step-up do
MVP ZTA.

Este backlog cobre criação da senha, armazenamento seguro, validação de PIN,
tentativas inválidas e bloqueio temporário.

## 2. Decisões fechadas

- A senha transacional é separada da senha de login.
- A senha transacional será um PIN numérico de 6 dígitos.
- A senha transacional deve ser armazenada apenas como hash.
- A senha transacional pode trafegar explicitamente apenas no endpoint de
  step-up, sempre sob HTTPS.
- A senha transacional não deve ser enviada pelo mobile como hash.
- A senha recebida não deve ser logada.
- Troca e reset de senha transacional ficam fora do MVP.

## 3. Criação da senha

Endpoint definido:

```http
POST /security/transaction-password
Authorization: Bearer <access_token>
```

Payload definido:

```json
{
  "transaction_password": "123456",
  "transaction_password_confirmation": "123456"
}
```

Regras:

- exige JWT válido;
- exige usuário ativo;
- é permitida para qualquer usuário ativo autenticado;
- não exige step-up token;
- não exige senha transacional anterior;
- exige confirmação no payload;
- bloqueia a criação se já existir senha transacional ativa;
- retorna sucesso sem expor senha, hash ou material sensível.

Fluxo:

```text
1. Usuário autenticado chama POST /security/transaction-password.
2. Backend valida JWT e identifica o usuário.
3. Backend valida se o usuário está ativo.
4. Backend valida PIN numérico de 6 dígitos.
5. Backend valida confirmação do PIN.
6. Backend verifica se já existe senha transacional ativa.
7. Backend gera hash.
8. Backend persiste a credencial.
9. Backend responde usando o envelope padrão da API.
```

## 4. Tentativas e bloqueio

Regras:

- o sistema conta tentativas inválidas consecutivas;
- após 3 tentativas inválidas, a senha transacional é bloqueada
  temporariamente;
- o bloqueio dura 30 minutos;
- após validação correta, o contador de falhas é zerado;
- o desbloqueio ocorre automaticamente por regra de leitura;
- quando o bloqueio expira, o contador de falhas é zerado;
- bloqueios mais fortes ficam para evolução futura baseada em risco.

## 5. Modelo conceitual

```text
transaction_passwords
- id
- user_id
- password_hash
- status
- failed_attempts
- locked_until
- created_at
- updated_at
- changed_at
```

Status iniciais:

```text
active
blocked
```

Mesmo com `status`, `locked_until` deve existir para representar o bloqueio
temporário de 30 minutos.

## 6. Erros cobertos

- `TRANSACTION_PASSWORD_ALREADY_SET`
- `TRANSACTION_PASSWORD_NOT_SET`
- `TRANSACTION_PASSWORD_INVALID`
- `TRANSACTION_PASSWORD_LOCKED`

O contrato completo de erros fica no backlog
[006d - zta-contracts-and-docs.md](<006d - zta-contracts-and-docs.md>).
