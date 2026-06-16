# Backlog API 016: Installation Identity Registration Flow

## 1. Status

- Tipo: Planning
- Área: Security
- Prioridade: High
- Estado: Discussão

## 2. Objetivo

Implementar o fluxo de registro explícito de uma nova instalação após login com
autorização restrita e step-up.

## 3. Escopo

- Permitir `POST /security/step-up/authorize` com contexto restrito.
- Reconhecer `POST /security/installations` como operação elegível para
  step-up.
- Exigir senha transacional ativa.
- Bloquear o fluxo se a senha estiver `not_set` ou `locked`.
- Implementar `POST /security/installations`.
- Validar autorização restrita, step-up e `X-Installation-Id`.
- Confirmar correspondência com a instalação apresentada no login.
- Confirmar vaga de forma atômica.
- Criar associação `known`.
- Consumir a autorização restrita.
- Criar sessão operacional vinculada.
- Retornar `access_token` e `refresh_token` normais.

## 4. Fora de escopo

- Bootstrap automático da primeira instalação no login.
- Listagem de instalações.
- Revogação de instalações.
- Recovery flow quando todas as instalações anteriores estiverem indisponíveis.

## 5. Dependências

- Backlog 013: domínio e repositórios.
- Backlog 014: sessão, JWT e contexto restrito.
- Backlog 015: login emitindo autorização restrita para nova instalação.

## 6. Orientação para tasks

As tasks deste backlog devem preservar a separação entre senha transacional,
step-up e registro da instalação. A senha transacional autoriza a operação; ela
não cria a instalação diretamente.

## 7. Referências

- [Installation Identity MVP](<010 - installation-identity-mvp.md>)
