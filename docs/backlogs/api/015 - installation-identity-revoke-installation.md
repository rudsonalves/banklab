# Backlog API 015: Revoke installation on `DELETE /security/installations/{installation_resource_id}`

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Revogar logicamente uma instalação cadastrada, preservando histórico e
paralisando o acesso revogado imediatamente.

## 3. Contrato

```http
DELETE /security/installations/{installation_resource_id}
Authorization: Bearer <access_token>
X-Installation-Id: <UUID v4>
```

## 4. Regras consolidadas

- não exige `step-up` no MVP;
- só pode revogar outra instalação do mesmo usuário;
- não pode remover a instalação vinculada à sessão atual;
- a revogação é lógica, com `status = revoked` e `revoked_at`;
- instalações revogadas permanecem no histórico;
- instalação revogada não volta a `known` por novo login.

## 5. Efeito operacional

Revogar uma instalação paralisa esse acesso imediatamente:

- todos os `access_token` associados à instalação revogada deixam de valer na
  hora;
- todos os `refresh_token` associados à instalação revogada devem ser
  invalidados no mesmo momento;
- a instalação revogada não pode concluir novas operações nem renovar sessão.

## 6. Tasks derivadas

- revogar logicamente a instalação;
- impedir autorrevogação da instalação em uso;
- invalidar sessão e tokens vinculados;
- retornar erro explícito quando a instalação atual for alvo da revogação.

## 7. Referências

- [Installation Identity MVP umbrella](<010 - installation-identity-mvp.md>)
- [List installations](<014 - installation-identity-list-installations.md>)
